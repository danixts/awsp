package cloudwatch

import "github.com/danixts/awsp/internal/gatewaytui"

var themes = gatewaytui.Themes

func themeIndexByName(name string) int {
	for i, t := range themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}
