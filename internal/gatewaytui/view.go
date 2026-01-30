package gatewaytui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}
	t := &Themes[m.themeIndex]
	currentTheme = t

	usableW := m.width * 94 / 100
	if usableW < panelMinWidth*2 {
		usableW = panelMinWidth * 2
	}

	half := usableW / 2
	panelW := half
	if panelW < panelMinWidth {
		panelW = panelMinWidth
	}
	panelH := m.height - 4
	if panelH < panelMinHeight {
		panelH = panelMinHeight
	}

	apiContent := m.apiList.View()
	if len(m.apis) == 0 && m.loading == "" && m.err == "" {
		apiContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render("Loading API Gateways...")
	} else if len(m.apis) == 0 && m.loading == "" {
		apiContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("No APIs (press R to retry)")
	} else if m.loading == "Loading APIs..." {
		apiContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Loading API Gateways...")
	}

	resourceContent := m.resourceList.View()
	if m.loading != "" && (m.loading == "Loading endpoints..." || m.loading == "Loading resources...") {
		resourceContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Loading...")
	} else if m.selectedAPI == nil {
		resourceContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("Select an API (left), then Enter")
	} else if len(m.endpoints) == 0 {
		resourceContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("No endpoints with methods")
	}

	profileStr := m.profile
	if profileStr == "" {
		profileStr = "default"
	}
	regionStr := m.region
	if regionStr == "" {
		regionStr = "(not set)"
	}
	apiTitleColor := t.InactiveBorder
	if m.focus == focusAPIs {
		apiTitleColor = t.ActiveBorder
	}
	apiTitleText := lipgloss.NewStyle().Foreground(lipgloss.Color(apiTitleColor)).Bold(true).Render(" API Gateways ")
	profileBadge := lipgloss.NewStyle().
		Background(lipgloss.Color(t.ProfileBg)).
		Foreground(lipgloss.Color(t.ProfileFg)).
		Bold(true).
		Padding(0, 1).
		Render(" PROFILE: " + profileStr + " ")
	regionBadge := lipgloss.NewStyle().
		Background(lipgloss.Color(t.ProfileBg)).
		Foreground(lipgloss.Color(t.ProfileFg)).
		Bold(true).
		Padding(0, 1).
		Render(" REGION: " + regionStr + " ")
	apiTitle := lipgloss.JoinHorizontal(lipgloss.Center, apiTitleText, " ", profileBadge, " ", regionBadge)

	activeBorder, inactiveBorder := PanelBorderStyle(t.ActiveBorder, t.InactiveBorder)
	endpointsTitleColor := t.InactiveBorder
	if m.focus == focusResources {
		endpointsTitleColor = t.ActiveBorder
	}
	endpointsTitle := lipgloss.NewStyle().Foreground(lipgloss.Color(endpointsTitleColor)).Bold(true).Render(" Endpoints ")
	apiBox := panelBorder(activeBorder, inactiveBorder, m.focus == focusAPIs, apiTitle, apiContent, panelW, panelH)
	resourceBox := panelBorder(activeBorder, inactiveBorder, m.focus == focusResources, endpointsTitle, resourceContent, panelW, panelH)

	row := lipgloss.JoinHorizontal(lipgloss.Top, apiBox, resourceBox)

	headerW := usableW
	if headerW < 50 {
		headerW = 50
	}
	profilePanelStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.ProfileBorder)).
		Foreground(lipgloss.Color(t.ProfileFg)).
		Background(lipgloss.Color(t.ProfileBg)).
		Bold(true).
		Padding(0, 1).
		Width(headerW).
		Align(lipgloss.Left)
	themeName := Themes[m.themeIndex].Name
	profilePanelContent := " Profile: " + profileStr + "  ·  Region: " + regionStr + "  ·  Theme: " + themeName + " "
	header := profilePanelStyle.Render(profilePanelContent)

	var body strings.Builder
	if m.loading != "" {
		body.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render("  " + m.spinner.View() + " " + m.loading))
		body.WriteString("\n")
	}
	if m.err != "" {
		body.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(t.Error)).Render("  " + m.err))
		body.WriteString("\n")
	}
	body.WriteString(row)
	body.WriteString("\n")
	help := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("  Tab/←/→ panels · Enter=select · R reload · t/T=theme · Esc back · q quit")
	body.WriteString(help)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, body.String())
	b := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, mainView)

	if m.modalState != modalNone {
		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.ModalBorder)).
			Padding(0, 1)
		if m.modalState == modalTimeSelect {
			modalW := usableW - 8
			if modalW > 70 {
				modalW = 70
			}
			modalH := m.height - 6
			if modalH > 22 {
				modalH = 22
			}
			timeTitle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Render(" Time range ") + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("· Esc cancel")
			box := modalStyle.Width(modalW).Height(modalH).Render(timeTitle + "\n\n" + m.modalTimeList.View())
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
		}
		logVp := m.modalLogView
		logVp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(t.LogContent))
		logContent := logVp.View()
		if m.loadingForModal && m.modalContent == "" {
			logContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Streaming...")
		}

		endpointInfo := ""
		if m.currentMethod != "" && m.currentPath != "" {
			methodStyle := themeMethodStyle(t, m.currentMethod)
			endpointInfo = " " + methodStyle.Render(" "+m.currentMethod+" ") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Render(m.currentPath) + " · "
		}
		logsTitle := endpointInfo + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Render("Logs") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("· Esc close")

		fullModal := modalStyle.
			Width(usableW).
			Height(m.height - 2).
			Render(" " + logsTitle + "\n\n" + logContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fullModal)
	}
	return b
}

func panelBorder(activeBorder, inactiveBorder lipgloss.Style, active bool, title string, content string, width, height int) string {
	style := inactiveBorder
	if active {
		style = activeBorder
	}
	body := title + "\n" + content
	return style.Width(width).Height(height).Render(body)
}
