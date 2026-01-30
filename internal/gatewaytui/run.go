package gatewaytui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Run(profile, region string) error {
	m := newModel(profile, region)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
