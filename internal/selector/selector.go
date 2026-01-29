package selector

import (
	"fmt"
	"strings"

	"github.com/danixts/awsp/internal/aws"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func RunProfileSelector(profiles []aws.Profile, current string) (string, error) {
	cyan := color.New(color.FgCyan).SprintFunc()

	tpl := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "► {{ .Name | cyan }}{{ if .Region }}  ·  {{ .Region | yellow }}{{ end }}",
		Inactive: "  {{ .Name }}{{ if .Region }}  ({{ .Region | faint }}){{ end }}",
		Selected: "✓ Selected: {{ .Name | green }}{{ if .Region }}  →  {{ .Region }}{{ end }}",
		Details:  "\n  Profile: {{ .Name }}{{ if .Region }}  |  Region: {{ .Region }}{{ end }}",
	}

	search := func(input string, i int) bool {
		p := profiles[i]
		name := strings.ReplaceAll(strings.ToLower(p.Name), " ", "")
		q := strings.ReplaceAll(strings.ToLower(input), " ", "")
		return strings.Contains(name, q)
	}

	sel := promptui.Select{
		Label:     fmt.Sprintf("Current profile: %s", cyan(current)),
		Items:     profiles,
		Templates: tpl,
		Size:      10,
		Searcher:  search,
	}

	idx, _, err := sel.Run()
	if err != nil {
		return "", err
	}

	return profiles[idx].Name, nil
}
