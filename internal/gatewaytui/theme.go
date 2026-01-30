package gatewaytui

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name string

	GETBg     string
	GETFg     string
	POSTBg    string
	POSTFg    string
	PUTBg     string
	PUTFg     string
	DELETEBg  string
	DELETEFg  string
	PATCHBg   string
	PATCHFg   string
	OPTIONSBg string
	OPTIONSFg string
	HEADBg    string
	HEADFg    string
	DefaultBg string
	DefaultFg string

	Accent        string
	ProfileBg     string
	ProfileBorder string
	ProfileFg     string
	ActiveBorder  string
	InactiveBorder string
	ModalBorder   string
	Loading       string
	Error         string
	Help          string
}

var themes = []Theme{
	{
		Name: "Purple",
		GETBg: "#22c55e", GETFg: "#fff",
		POSTBg: "#3b82f6", POSTFg: "#fff",
		PUTBg: "#f97316", PUTFg: "#fff",
		DELETEBg: "#ef4444", DELETEFg: "#fff",
		PATCHBg: "#eab308", PATCHFg: "#1f2937",
		OPTIONSBg: "#6b7280", OPTIONSFg: "#fff",
		HEADBg: "#8b5cf6", HEADFg: "#fff",
		DefaultBg: "#64748b", DefaultFg: "#fff",
		Accent: "#7D56F4", ProfileBg: "#5b21b6", ProfileBorder: "#7D56F4", ProfileFg: "#fff",
		ActiveBorder: "#7D56F4", InactiveBorder: "#3D3D3D", ModalBorder: "#7D56F4",
		Loading: "#7D56F4", Error: "#FF6B6B", Help: "#6B6B6B",
	},
	{
		Name: "Ocean",
		GETBg: "#0ea5e9", GETFg: "#fff",
		POSTBg: "#06b6d4", POSTFg: "#fff",
		PUTBg: "#0284c7", PUTFg: "#fff",
		DELETEBg: "#dc2626", DELETEFg: "#fff",
		PATCHBg: "#38bdf8", PATCHFg: "#0c4a6e",
		OPTIONSBg: "#64748b", OPTIONSFg: "#fff",
		HEADBg: "#0369a1", HEADFg: "#fff",
		DefaultBg: "#475569", DefaultFg: "#fff",
		Accent: "#0ea5e9", ProfileBg: "#0c4a6e", ProfileBorder: "#0ea5e9", ProfileFg: "#e0f2fe",
		ActiveBorder: "#0ea5e9", InactiveBorder: "#334155", ModalBorder: "#0ea5e9",
		Loading: "#0ea5e9", Error: "#f87171", Help: "#94a3b8",
	},
	{
		Name: "Forest",
		GETBg: "#22c55e", GETFg: "#fff",
		POSTBg: "#16a34a", POSTFg: "#fff",
		PUTBg: "#ca8a04", PUTFg: "#fff",
		DELETEBg: "#b91c1c", DELETEFg: "#fff",
		PATCHBg: "#84cc16", PATCHFg: "#14532d",
		OPTIONSBg: "#4b5563", OPTIONSFg: "#fff",
		HEADBg: "#15803d", HEADFg: "#fff",
		DefaultBg: "#404040", DefaultFg: "#fff",
		Accent: "#22c55e", ProfileBg: "#14532d", ProfileBorder: "#22c55e", ProfileFg: "#dcfce7",
		ActiveBorder: "#22c55e", InactiveBorder: "#374151", ModalBorder: "#22c55e",
		Loading: "#22c55e", Error: "#f87171", Help: "#6b7280",
	},
	{
		Name: "Sunset",
		GETBg: "#f97316", GETFg: "#fff",
		POSTBg: "#ea580c", POSTFg: "#fff",
		PUTBg: "#dc2626", PUTFg: "#fff",
		DELETEBg: "#991b1b", DELETEFg: "#fff",
		PATCHBg: "#fbbf24", PATCHFg: "#78350f",
		OPTIONSBg: "#78716c", OPTIONSFg: "#fff",
		HEADBg: "#c2410c", HEADFg: "#fff",
		DefaultBg: "#57534e", DefaultFg: "#fff",
		Accent: "#f97316", ProfileBg: "#9a3412", ProfileBorder: "#fb923c", ProfileFg: "#ffedd5",
		ActiveBorder: "#f97316", InactiveBorder: "#44403c", ModalBorder: "#f97316",
		Loading: "#f97316", Error: "#ef4444", Help: "#a8a29e",
	},
	{
		Name: "Neon",
		GETBg: "#00ff88", GETFg: "#000",
		POSTBg: "#00d4ff", POSTFg: "#000",
		PUTBg: "#ff00ff", PUTFg: "#fff",
		DELETEBg: "#ff0040", DELETEFg: "#fff",
		PATCHBg: "#ffff00", PATCHFg: "#000",
		OPTIONSBg: "#00ffff", OPTIONSFg: "#000",
		HEADBg: "#bf00ff", HEADFg: "#fff",
		DefaultBg: "#808080", DefaultFg: "#000",
		Accent: "#00ff88", ProfileBg: "#0d0d0d", ProfileBorder: "#00ff88", ProfileFg: "#00ff88",
		ActiveBorder: "#00ff88", InactiveBorder: "#333333", ModalBorder: "#00ff88",
		Loading: "#00ff88", Error: "#ff0040", Help: "#666666",
	},
	{
		Name: "Cyber",
		GETBg: "#06b6d4", GETFg: "#000",
		POSTBg: "#8b5cf6", POSTFg: "#fff",
		PUTBg: "#f59e0b", PUTFg: "#000",
		DELETEBg: "#ef4444", DELETEFg: "#fff",
		PATCHBg: "#10b981", PATCHFg: "#fff",
		OPTIONSBg: "#6366f1", OPTIONSFg: "#fff",
		HEADBg: "#ec4899", HEADFg: "#fff",
		DefaultBg: "#475569", DefaultFg: "#fff",
		Accent: "#06b6d4", ProfileBg: "#0f172a", ProfileBorder: "#06b6d4", ProfileFg: "#22d3ee",
		ActiveBorder: "#06b6d4", InactiveBorder: "#1e293b", ModalBorder: "#06b6d4",
		Loading: "#06b6d4", Error: "#f43f5e", Help: "#64748b",
	},
}

var currentTheme *Theme

func themeMethodStyle(t *Theme, method string) lipgloss.Style {
	var bg, fg string
	switch method {
	case "GET":
		bg, fg = t.GETBg, t.GETFg
	case "POST":
		bg, fg = t.POSTBg, t.POSTFg
	case "PUT":
		bg, fg = t.PUTBg, t.PUTFg
	case "DELETE":
		bg, fg = t.DELETEBg, t.DELETEFg
	case "PATCH":
		bg, fg = t.PATCHBg, t.PATCHFg
	case "OPTIONS":
		bg, fg = t.OPTIONSBg, t.OPTIONSFg
	case "HEAD":
		bg, fg = t.HEADBg, t.HEADFg
	default:
		bg, fg = t.DefaultBg, t.DefaultFg
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg)).
		Bold(true).
		Width(methodWidth).
		Align(lipgloss.Center)
}

func themeIndexByName(name string) int {
	for i, t := range themes {
		if t.Name == name {
			return i
		}
	}
	return 0
}
