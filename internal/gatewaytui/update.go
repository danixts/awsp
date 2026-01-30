package gatewaytui

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danixts/awsp/internal/apigateway"
	"github.com/danixts/awsp/internal/store"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.loading != "" || m.loadingForModal {
			return m, cmd
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		usableW := msg.Width * 94 / 100
		if usableW < panelMinWidth*2 {
			usableW = panelMinWidth * 2
		}
		half := usableW / 2
		panelW := half
		if panelW < panelMinWidth {
			panelW = panelMinWidth
		}
		contentH := msg.Height - 4
		if contentH < 6 {
			contentH = 6
		}
		contentW := panelW - 4
		if contentW < 12 {
			contentW = 12
		}
		m.apiList.SetSize(contentW, contentH)
		m.resourceList.SetSize(contentW, contentH)
		modalTimeW := usableW - 8
		if modalTimeW > 70 {
			modalTimeW = 70
		}
		modalTimeH := msg.Height - 6
		if modalTimeH > 22 {
			modalTimeH = 22
		}
		m.modalTimeList.SetSize(modalTimeW-4, modalTimeH-4)
		m.modalLogView.Width = usableW - 4
		m.modalLogView.Height = msg.Height - 8
		return m, nil

	case apisLoadedMsg:
		m.loading = ""
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.apis = msg.apis
		m.cachedAPIs = msg.apis
		if len(msg.apis) == 0 {
			m.err = "No API Gateways found"
		}
		items := apiItemsFrom(msg.apis)
		cmds := []tea.Cmd{m.apiList.SetItems(items)}
		for _, a := range msg.apis {
			cmds = append(cmds, loadResourcesCmd(a.ID, true))
		}
		return m, tea.Batch(cmds...)

	case resourcesLoadedMsg:
		if msg.preload {
			if msg.err == nil {
				m.cachedResources[msg.apiID] = msg.resources
			}
			return m, nil
		}
		m.loading = ""
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		for i := range m.apis {
			if m.apis[i].ID == msg.apiID {
				m.selectedAPI = &m.apis[i]
				break
			}
		}
		m.endpoints = apigateway.EndpointsFromResources(msg.resources)
		items := endpointItemsFrom(m.endpoints)
		c := m.resourceList.SetItems(items)
		m.resourceList.ResetSelected()
		m.focus = focusResources
		return m, c

	case logGroupLoadedMsg:
		m.loading = ""
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.cachedLogGroups[msg.logKey] = msg.logGroup
		m.pendingLogGroup = msg.logGroup
		m.modalState = modalTimeSelect
		items := timeRangeItemsFrom(apigateway.DefaultTimeRanges)
		c := m.modalTimeList.SetItems(items)
		return m, c

	case logStreamStartedMsg:
		m.loadingForModal = false
		m.logStreamCh = msg.ch
		m.modalContent = ""
		m.modalLogLines = nil
		m.modalLogView.SetContent("")
		m.modalState = modalLogView
		return m, waitForLogLineCmd(msg.ch)

	case logLineMsg:
		if m.modalState == modalLogView {
			m.modalLogLines = append(m.modalLogLines, msg.line)
			if len(m.modalLogLines) > maxLogLinesTail {
				m.modalLogLines = m.modalLogLines[len(m.modalLogLines)-maxLogLinesTail:]
			}
			m.modalContent = strings.Join(m.modalLogLines, "\n")
			m.modalLogView.SetContent(m.modalContent)
			m.modalLogView.GotoBottom()
			return m, waitForLogLineCmd(m.logStreamCh)
		}
		return m, nil

	case logStreamDoneMsg:
		return m, nil

	case logsLoadedMsg:
		m.loading = ""
		if m.loadingForModal {
			return m, nil
		}
		m.logViewport.SetContent(msg.content)
		if msg.err != nil {
			m.err = msg.err.Error()
		}
		return m, nil

	case tea.KeyMsg:
		if m.modalState != modalNone {
			if msg.String() == "esc" {
				if m.modalState == modalLogView {
					if m.cancelLogStream != nil {
						m.cancelLogStream()
						m.cancelLogStream = nil
					}
					m.focus = focusResources
				}
				m.modalState = modalNone
				m.pendingLogGroup = ""
				m.currentMethod = ""
				m.currentPath = ""
				return m, nil
			}
			if m.modalState == modalTimeSelect {
				if msg.String() == "enter" {
					it := m.modalTimeList.SelectedItem()
					if t, ok := it.(timeRangeItem); ok {
						m.loadingForModal = true
						m.modalState = modalLogView
						m.modalContent = ""
						m.modalLogLines = nil
						ctx, cancel := context.WithCancel(context.Background())
						m.cancelLogStream = cancel
						return m, tea.Batch(loadLogsStreamCmd(ctx, m.pendingLogGroup, t.Since, m.region), func() tea.Msg { return m.spinner.Tick() })
					}
				}
				var cmd tea.Cmd
				m.modalTimeList, cmd = m.modalTimeList.Update(msg)
				return m, cmd
			}
			if m.modalState == modalLogView {
				var cmd tea.Cmd
				m.modalLogView, cmd = m.modalLogView.Update(msg)
				return m, cmd
			}
		}
		switch msg.String() {
		case "ctrl+c", "q":
			if m.cancelLogStream != nil {
				m.cancelLogStream()
			}
			m.quitting = true
			return m, tea.Quit
		case "t":
			if m.modalState == modalNone {
				m.themeIndex = (m.themeIndex + 1) % len(Themes)
				_ = store.SaveGatewayTheme(Themes[m.themeIndex].Name)
				t := &Themes[m.themeIndex]
				m.apiList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				m.resourceList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				return m, nil
			}
		case "T":
			if m.modalState == modalNone {
				m.themeIndex = (m.themeIndex - 1 + len(Themes)) % len(Themes)
				_ = store.SaveGatewayTheme(Themes[m.themeIndex].Name)
				t := &Themes[m.themeIndex]
				m.apiList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				m.resourceList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				return m, nil
			}
		case "tab", "right", "shift+tab", "left":
			newFocus := (m.focus + 1) % 2
			if newFocus == focusAPIs {
				m.endpoints = nil
				m.selectedAPI = nil
				m.resourceList.SetItems([]list.Item{})
			}
			m.focus = newFocus
			return m, nil
		case "r", "R":
			if m.focus == focusAPIs {
				m.loading = "Loading APIs..."
				return m, tea.Batch(loadAPIsCmd(), func() tea.Msg { return m.spinner.Tick() })
			}
			if m.focus == focusResources && m.selectedAPI != nil {
				delete(m.cachedResources, m.selectedAPI.ID)
				m.loading = "Loading resources..."
				return m, tea.Batch(loadResourcesCmd(m.selectedAPI.ID, false), func() tea.Msg { return m.spinner.Tick() })
			}
		case "enter":
			if m.focus == focusAPIs {
				it := m.apiList.SelectedItem()
				if it == nil {
					return m, nil
				}
				if a, ok := it.(apiItem); ok {
					m.err = ""
					if res, ok := m.cachedResources[a.ID]; ok {
						for i := range m.apis {
							if m.apis[i].ID == a.ID {
								m.selectedAPI = &m.apis[i]
								break
							}
						}
						m.endpoints = apigateway.EndpointsFromResources(res)
						items := endpointItemsFrom(m.endpoints)
						cmd := m.resourceList.SetItems(items)
						m.resourceList.ResetSelected()
						m.focus = focusResources
						return m, cmd
					}
					m.loading = "Loading endpoints..."
					return m, tea.Batch(loadResourcesCmd(a.ID, false), func() tea.Msg { return m.spinner.Tick() })
				}
			}
			if m.focus == focusResources && m.selectedAPI != nil {
				it := m.resourceList.SelectedItem()
				if it == nil {
					return m, nil
				}
				if ep, ok := it.(endpointItem); ok {
					m.currentMethod = ep.Method
					m.currentPath = ep.Path
					logKey := m.selectedAPI.ID + ":" + ep.ResourceID + ":" + ep.Method
					if logGroup, ok := m.cachedLogGroups[logKey]; ok {
						m.pendingLogGroup = logGroup
						m.modalState = modalTimeSelect
						items := timeRangeItemsFrom(apigateway.DefaultTimeRanges)
						c := m.modalTimeList.SetItems(items)
						return m, c
					}
					m.loading = "Loading log group..."
					return m, tea.Batch(
						loadLogGroupCmd(m.selectedAPI.ID, ep.ResourceID, ep.Method, logKey),
						func() tea.Msg { return m.spinner.Tick() },
					)
				}
			}
		}
	}

	var cmd tea.Cmd
	switch m.focus {
	case focusAPIs:
		m.apiList, cmd = m.apiList.Update(msg)
	case focusResources:
		m.resourceList, cmd = m.resourceList.Update(msg)
	}
	return m, cmd
}
