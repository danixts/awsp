package selector

import (
	"fmt"
	"strings"

	"github.com/danixts/awsp/internal/aws"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func RunRegionSelector(regions []aws.Region, current string) (string, error) {
	cyan := color.New(color.FgCyan).SprintFunc()

	tpl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "► {{ .Code | cyan }}  ·  {{ .Name }}",
		Inactive: "  {{ .Code }}  ({{ .Name | faint }})",
		Selected: "✓ Selected: {{ .Code | green }}  →  {{ .Name }}",
		Details:  "\n  Region: {{ .Code }}  |  {{ .Name }}",
	}

	search := func(input string, i int) bool {
		r := regions[i]
		text := strings.ReplaceAll(strings.ToLower(r.Code+" "+r.Name), " ", "")
		q := strings.ReplaceAll(strings.ToLower(input), " ", "")
		return strings.Contains(text, q)
	}

	label := "Select region:"
	if current != "" {
		label = fmt.Sprintf("Current region: %s", cyan(current))
	}

	sel := promptui.Select{
		Label:     label,
		Items:     regions,
		Templates: tpl,
		Size:      12,
		Searcher:  search,
	}

	idx, _, err := sel.Run()
	if err != nil {
		return "", err
	}

	return regions[idx].Code, nil
}
