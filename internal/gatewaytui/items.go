package gatewaytui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danixts/awsp/internal/apigateway"
)

type timeRangeItem struct {
	Label string
	Since string
}

func (t timeRangeItem) Title() string       { return t.Label }
func (t timeRangeItem) Description() string { return "" }
func (t timeRangeItem) FilterValue() string { return t.Label }

type apiItem struct {
	apigateway.RestAPI
}

func (a apiItem) Title() string       { return a.Name + " · " + a.ID }
func (a apiItem) Description() string { return "" }
func (a apiItem) FilterValue() string { return a.Name + " " + a.ID }

type endpointItem struct {
	apigateway.Endpoint
}

func (e endpointItem) Title() string       { return e.Method + "  " + e.Path }
func (e endpointItem) Description() string { return "" }
func (e endpointItem) FilterValue() string { return e.Method + " " + e.Path }

const methodWidth = 7

type apiDelegate struct {
	list.DefaultDelegate
}

func (d apiDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	api, ok := item.(apiItem)
	if !ok {
		d.DefaultDelegate.Render(w, m, index, item)
		return
	}
	t := currentTheme
	if t == nil {
		t = &Themes[0]
	}

	nameStyle := lipgloss.NewStyle()
	idStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	if index == m.Index() {
		nameStyle = nameStyle.Foreground(lipgloss.Color(t.Accent)).Bold(true)
		idStyle = idStyle.Foreground(lipgloss.Color(t.Accent))
	}

	line := nameStyle.Render(api.Name) + idStyle.Render(" · "+api.ID)
	fmt.Fprint(w, line)
}

func (d apiDelegate) Height() int   { return d.DefaultDelegate.Height() }
func (d apiDelegate) Spacing() int { return d.DefaultDelegate.Spacing() }
func (d apiDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

type endpointDelegate struct {
	list.DefaultDelegate
}

func (d endpointDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	ep, ok := item.(endpointItem)
	if !ok {
		d.DefaultDelegate.Render(w, m, index, item)
		return
	}
	method := strings.ToUpper(ep.Method)
	t := currentTheme
	if t == nil {
		t = &Themes[0]
	}
	style := themeMethodStyle(t, method)
	methodBadge := style.Render(method)
	pathStyle := lipgloss.NewStyle().PaddingLeft(1)
	if index == m.Index() {
		pathStyle = pathStyle.Foreground(lipgloss.Color(t.Accent)).Bold(true)
	}
	line := methodBadge + pathStyle.Render(ep.Path)
	fmt.Fprint(w, line)
}

func (d endpointDelegate) Height() int   { return d.DefaultDelegate.Height() }
func (d endpointDelegate) Spacing() int { return d.DefaultDelegate.Spacing() }
func (d endpointDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

type timeRangeDelegate struct {
	list.DefaultDelegate
}

func (d timeRangeDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	tr, ok := item.(timeRangeItem)
	if !ok {
		d.DefaultDelegate.Render(w, m, index, item)
		return
	}
	t := currentTheme
	if t == nil {
		t = &Themes[0]
	}

	style := lipgloss.NewStyle()
	if index == m.Index() {
		style = style.Foreground(lipgloss.Color(t.Accent)).Bold(true)
	}

	line := style.Render(tr.Label)
	fmt.Fprint(w, line)
}

func (d timeRangeDelegate) Height() int   { return d.DefaultDelegate.Height() }
func (d timeRangeDelegate) Spacing() int { return d.DefaultDelegate.Spacing() }
func (d timeRangeDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return d.DefaultDelegate.Update(msg, m)
}

func apiItemsFrom(apis []apigateway.RestAPI) []list.Item {
	out := make([]list.Item, 0, len(apis))
	for _, a := range apis {
		out = append(out, apiItem{a})
	}
	return out
}

func endpointItemsFrom(endpoints []apigateway.Endpoint) []list.Item {
	out := make([]list.Item, 0, len(endpoints))
	for _, e := range endpoints {
		out = append(out, endpointItem{e})
	}
	return out
}

func timeRangeItemsFrom(ranges []apigateway.TimeRange) []list.Item {
	out := make([]list.Item, 0, len(ranges))
	for _, r := range ranges {
		out = append(out, timeRangeItem{Label: r.Label, Since: r.Since})
	}
	return out
}
