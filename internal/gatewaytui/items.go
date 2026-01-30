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

type themeItem struct {
	name  string
	index int
}

func (t themeItem) Title() string       { return t.name }
func (t themeItem) Description() string { return "" }
func (t themeItem) FilterValue() string { return t.name }

func themeItemsFrom(themes []Theme) []list.Item {
	out := make([]list.Item, 0, len(themes))
	for i := range themes {
		out = append(out, themeItem{name: themes[i].Name, index: i})
	}
	return out
}

const methodWidth = 7

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
		t = &themes[0]
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
