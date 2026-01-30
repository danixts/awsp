package selector

import (
	"fmt"
	"strings"

	"github.com/danixts/awsp/internal/apigateway"
	"github.com/manifoldco/promptui"
)

type gatewayAPIItem struct {
	Label   string
	API     *apigateway.RestAPI
	IsReload bool
}

func RunGatewayAPISelector(apis []apigateway.RestAPI, withReload bool) (*apigateway.RestAPI, bool, error) {
	if len(apis) == 0 {
		return nil, false, fmt.Errorf("no APIs found")
	}
	items := make([]gatewayAPIItem, 0, len(apis)+1)
	for i := range apis {
		a := &apis[i]
		items = append(items, gatewayAPIItem{Label: a.Name + "  ·  " + a.ID, API: a, IsReload: false})
	}
	if withReload {
		items = append(items, gatewayAPIItem{Label: "🔄 Reload API list", API: nil, IsReload: true})
	}
	tpl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "► {{ .Label | cyan }}",
		Inactive: "  {{ .Label }}",
		Selected: "✓ {{ .Label | green }}",
		Details:  "{{ if .API }}\n  API: {{ .API.Name }}  |  ID: {{ .API.ID }}{{ end }}",
	}
	search := func(input string, i int) bool {
		it := items[i]
		if it.IsReload {
			return strings.Contains(strings.ToLower(it.Label), strings.ReplaceAll(strings.ToLower(input), " ", ""))
		}
		q := strings.ReplaceAll(strings.ToLower(input), " ", "")
		return strings.Contains(strings.ToLower(it.API.Name), q) || strings.Contains(strings.ToLower(it.API.ID), q)
	}
	sel := promptui.Select{
		Label:     "Select API Gateway:",
		Items:     items,
		Templates: tpl,
		Size:      10,
		Searcher:  search,
	}
	idx, _, err := sel.Run()
	if err != nil {
		return nil, false, err
	}
	it := items[idx]
	if it.IsReload {
		return nil, true, nil
	}
	return it.API, false, nil
}

const (
	viewLogsOption = "View logs for an endpoint"
	backOption     = "← Back to API list"
	reloadOption   = "🔄 Reload (refresh resources)"
	exitOption     = "Exit"
)

type AfterResourcesChoice int

const (
	ChoiceViewLogs AfterResourcesChoice = iota
	ChoiceBack
	ChoiceReload
	ChoiceExit
)

func RunAfterResourcesSelector() (AfterResourcesChoice, error) {
	items := []string{viewLogsOption, backOption, reloadOption, exitOption}
	sel := promptui.Select{
		Label: "What next?",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Active:   "► {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}
	idx, _, err := sel.Run()
	if err != nil {
		return 0, err
	}
	return AfterResourcesChoice(idx), nil
}

func RunEndpointSelector(endpoints []apigateway.Endpoint) (*apigateway.Endpoint, error) {
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints with methods")
	}
	tpl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "► {{ .Path | cyan }}  {{ .Method }}",
		Inactive: "  {{ .Path }}  {{ .Method }}",
		Selected: "✓ {{ .Path | green }}  {{ .Method }}",
		Details:  "\n  Path: {{ .Path }}  |  Method: {{ .Method }}",
	}
	search := func(input string, i int) bool {
		e := endpoints[i]
		q := strings.ReplaceAll(strings.ToLower(input), " ", "")
		return strings.Contains(strings.ToLower(e.Path), q) || strings.Contains(strings.ToLower(e.Method), q)
	}
	sel := promptui.Select{
		Label:     "Select endpoint:",
		Items:     endpoints,
		Templates: tpl,
		Size:      12,
		Searcher:  search,
	}
	idx, _, err := sel.Run()
	if err != nil {
		return nil, err
	}
	return &endpoints[idx], nil
}

func RunTimeRangeSelector(ranges []apigateway.TimeRange) (since string, err error) {
	if len(ranges) == 0 {
		return "", fmt.Errorf("no time ranges")
	}
	tpl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "► {{ .Label | cyan }}",
		Inactive: "  {{ .Label }}",
		Selected: "✓ {{ .Label | green }}",
	}
	sel := promptui.Select{
		Label:     "Time range:",
		Items:     ranges,
		Templates: tpl,
		Size:      10,
	}
	idx, _, err := sel.Run()
	if err != nil {
		return "", err
	}
	return ranges[idx].Since, nil
}
