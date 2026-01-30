package apigateway

import (
	"fmt"
	"strings"
)

func FormatResourcesTable(resources []ResourceWithMethods) string {
	if len(resources) == 0 {
		return "  (no resources with methods)\n"
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
	fmt.Fprintf(&b, "  %-*s  %s\n", pathW, "path", "methods")
	fmt.Fprintf(&b, "  %-*s  %s\n", pathW, "----", "-------")
	for _, r := range resources {
		methods := strings.Join(r.Methods, ", ")
		fmt.Fprintf(&b, "  %-*s  %s\n", pathW, r.Path, methods)
	}
	return b.String()
}
