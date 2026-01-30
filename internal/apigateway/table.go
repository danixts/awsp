package apigateway

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

var (
	tableHeader = color.New(color.FgHiBlack).SprintFunc()
	tableDim    = color.New(color.FgHiBlack).SprintFunc()
)

func FormatResourcesTable(resources []ResourceWithMethods) string {
	if len(resources) == 0 {
		return "  " + tableDim("(no resources with methods)") + "\n"
	}
	maxPath := 4
	for _, r := range resources {
		if len(r.Path) > maxPath {
			maxPath = len(r.Path)
		}
	}
	pathW := maxPath + 2
	if pathW < 12 {
		pathW = 12
	}
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(tableHeader(fmt.Sprintf("%-*s  %s", pathW, "path", "methods")))
	b.WriteString("\n  ")
	b.WriteString(tableDim(fmt.Sprintf("%-*s  %s", pathW, "──", "───────")))
	b.WriteString("\n")
	for _, r := range resources {
		methods := strings.Join(r.Methods, ", ")
		fmt.Fprintf(&b, "  %-*s  %s\n", pathW, r.Path, methods)
	}
	return b.String()
}
