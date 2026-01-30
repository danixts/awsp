package cloudwatch

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/danixts/awsp/internal/gatewaytui"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}
	t := &themes[m.themeIndex]

	usableW := m.width * 94 / 100
	if usableW < panelMinWidth {
		usableW = panelMinWidth
	}

	panelW := usableW - 2
	if panelW < panelMinWidth {
		panelW = panelMinWidth
	}
	panelH := m.height - 4
	if panelH < panelMinHeight {
		panelH = panelMinHeight
	}

	listContent := m.logGroupList.View()
	if len(m.logGroups) == 0 && m.loading == "" && m.err == "" {
		listContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render("Loading log groups...")
	} else if len(m.logGroups) == 0 && m.loading == "" {
		listContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("No log groups (press R to retry)")
	} else if m.loading == "Loading log groups..." {
		listContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Loading log groups...")
	}

	activeBorder, _ := panelBorderStyle(t.ActiveBorder, t.InactiveBorder)
	mainBox := activeBorder.Width(panelW).Height(panelH).Render(" Log Groups \n" + listContent)

	profileStr := m.profile
	if profileStr == "" {
		profileStr = "default"
	}
	regionStr := m.region
	if regionStr == "" {
		regionStr = "(not set)"
	}
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
	themeName := themes[m.themeIndex].Name
	profilePanelContent := " CloudWatch Logs  ·  Profile: " + profileStr + "  ·  Region: " + regionStr + "  ·  Theme: " + themeName + " "
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
	body.WriteString(mainBox)
	body.WriteString("\n")
	help := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("  Enter=stream logs · R reload · t/T=theme · / filter · Esc back · q quit")
	body.WriteString(help)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, body.String())
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, mainView)

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

		logGroupInfo := ""
		if m.currentLogGroup != "" {
			logGroupStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(t.Accent)).
				Bold(true)
			logGroupInfo = logGroupStyle.Render(m.currentLogGroup) + " · "
		}
		logsTitle := logGroupInfo + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Accent)).Render("Logs") + " " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("· Esc close")

		fullModal := modalStyle.
			Width(usableW).
			Height(m.height - 2).
			Render(" " + logsTitle + "\n\n" + logContent)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, fullModal)
	}
	return centered
}

func panelBorderStyle(activeBorderColor, inactiveBorderColor string) (active, inactive lipgloss.Style) {
	return gatewaytui.PanelBorderStyle(activeBorderColor, inactiveBorderColor)
}
