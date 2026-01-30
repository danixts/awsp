package gatewaytui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}
	t := &themes[m.themeIndex]
	currentTheme = t

	half := m.width / 2
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
		apiContent = "\n  Loading API Gateways..."
	} else if len(m.apis) == 0 && m.loading == "" {
		apiContent = "\n  No APIs (press R to retry)"
	} else if m.loading == "Loading APIs..." {
		apiContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Loading API Gateways...")
	}

	resourceContent := m.resourceList.View()
	if m.loading != "" && (m.loading == "Loading endpoints..." || m.loading == "Loading resources...") {
		resourceContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Loading...")
	} else if m.selectedAPI == nil {
		resourceContent = "\n  Select an API (left), then Enter"
	} else if len(m.endpoints) == 0 {
		resourceContent = "\n  No endpoints with methods"
	}

	activeBorder, inactiveBorder := panelBorderStyle(t.ActiveBorder, t.InactiveBorder)
	apiBox := panelBorder(activeBorder, inactiveBorder, m.focus == focusAPIs, " API Gateways ", apiContent, panelW, panelH)
	resourceBox := panelBorder(activeBorder, inactiveBorder, m.focus == focusResources, " Endpoints ", resourceContent, panelW, panelH)

	row := lipgloss.JoinHorizontal(lipgloss.Top, apiBox, resourceBox)

	profileStr := m.profile
	if profileStr == "" {
		profileStr = "default"
	}
	regionStr := m.region
	if regionStr == "" {
		regionStr = "(not set)"
	}
	headerW := m.width
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
	profilePanelContent := " Profile: " + profileStr + "  ·  Region: " + regionStr + "  ·  t=theme "
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
	help := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Help)).Render("  API Gateways (left) · Endpoints (right) · Enter=load/open logs · Tab/←/→ · R reload · t=theme · Esc · q quit")
	body.WriteString(help)

	mainView := lipgloss.JoinVertical(lipgloss.Left, header, body.String())
	b := mainView

	if m.modalState != modalNone {
		modalStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.ModalBorder)).
			Padding(0, 1)
		if m.modalState == modalThemeSelect {
			modalW := m.width - 12
			if modalW > 44 {
				modalW = 44
			}
			modalH := m.height - 8
			if modalH > 18 {
				modalH = 18
			}
			box := modalStyle.Width(modalW).Height(modalH).Render(" Theme · Enter=apply · Esc cancel\n\n" + m.modalThemeList.View())
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
		}
		if m.modalState == modalTimeSelect {
			modalW := m.width - 8
			if modalW > 70 {
				modalW = 70
			}
			modalH := m.height - 6
			if modalH > 22 {
				modalH = 22
			}
			box := modalStyle.Width(modalW).Height(modalH).Render(" Time range · Esc cancel\n\n" + m.modalTimeList.View())
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
		}
		logContent := m.modalLogView.View()
		if m.loadingForModal && m.modalContent == "" {
			logContent = "\n  " + lipgloss.NewStyle().Foreground(lipgloss.Color(t.Loading)).Render(m.spinner.View()+" Streaming...")
		}
		fullModal := modalStyle.
			Width(m.width).
			Height(m.height).
			Render(" Logs · Esc close\n\n" + logContent)
		return fullModal
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
