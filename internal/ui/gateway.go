package ui

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	cyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	dim     = color.New(color.FgHiBlack).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
)

func PrintGatewayHeader() {
	fmt.Println()
	fmt.Println(dim("┌────────────────────────────────────────┐"))
	fmt.Println(dim("│  ") + cyan("  API Gateway") + dim("                         │"))
	fmt.Println(dim("│  ") + dim("  List · Endpoints · Logs") + dim("               │"))
	fmt.Println(dim("└────────────────────────────────────────┘"))
	fmt.Println()
}

func PrintSection(title string) {
	fmt.Println(dim("  ── ") + magenta(title) + dim(" ──"))
	fmt.Println()
}

func PrintAPISelected(name string) {
	fmt.Println(dim("  › ") + green(name))
	fmt.Println()
}

func PrintWhatNext() {
	fmt.Println(dim("  ── ") + yellow("What next?") + dim(" ──"))
	fmt.Println()
}
