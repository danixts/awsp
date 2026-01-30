package cloudwatch

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danixts/awsp/internal/domain"
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
		if usableW < panelMinWidth {
			usableW = panelMinWidth
		}
		contentH := msg.Height - 4
		if contentH < 6 {
			contentH = 6
		}
		contentW := usableW - 4
		if contentW < 12 {
			contentW = 12
		}
		m.logGroupList.SetSize(contentW, contentH)
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

	case logGroupsLoadedMsg:
		m.loading = ""
		if msg.err != nil {
			m.err = msg.err.Error()
			return m, nil
		}
		m.err = ""
		m.logGroups = msg.groups
		if len(msg.groups) == 0 {
			m.err = "No log groups found"
		}
		items := logGroupItemsFrom(msg.groups)
		cmd := m.logGroupList.SetItems(items)
		return m, cmd

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

	case tea.KeyMsg:
		if m.modalState != modalNone {
			if msg.String() == "esc" {
				if m.modalState == modalLogView {
					if m.cancelLogStream != nil {
						m.cancelLogStream()
						m.cancelLogStream = nil
					}
				}
				m.modalState = modalNone
				m.pendingLogGroup = ""
				m.currentLogGroup = ""
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
						return m, tea.Batch(
							loadLogsStreamCmd(ctx, m.pendingLogGroup, t.Since, m.region),
							func() tea.Msg { return m.spinner.Tick() },
						)
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
				m.themeIndex = (m.themeIndex + 1) % len(themes)
				_ = store.SaveGatewayTheme(themes[m.themeIndex].Name)
				t := &themes[m.themeIndex]
				m.logGroupList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				return m, nil
			}
		case "T":
			if m.modalState == modalNone {
				m.themeIndex = (m.themeIndex - 1 + len(themes)) % len(themes)
				_ = store.SaveGatewayTheme(themes[m.themeIndex].Name)
				t := &themes[m.themeIndex]
				m.logGroupList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color(t.ActiveBorder)).Bold(true)
				return m, nil
			}
		case "r", "R":
			m.loading = "Loading log groups..."
			return m, tea.Batch(loadLogGroupsCmd(""), func() tea.Msg { return m.spinner.Tick() })
		case "enter":
			it := m.logGroupList.SelectedItem()
			if it == nil {
				return m, nil
			}
			if g, ok := it.(logGroupItem); ok {
				m.err = ""
				m.selectedGroup = &g.LogGroup
				m.pendingLogGroup = g.Name
				m.currentLogGroup = g.Name
				m.modalState = modalTimeSelect
				items := timeRangeItemsFrom(domain.DefaultTimeRanges)
				c := m.modalTimeList.SetItems(items)
				return m, c
			}
		}
	}

	var cmd tea.Cmd
	m.logGroupList, cmd = m.logGroupList.Update(msg)
	return m, cmd
}
