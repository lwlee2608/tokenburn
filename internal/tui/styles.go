package tui

import "github.com/charmbracelet/lipgloss"

const barPadding = 30 // space taken by "  Used: XXX%  " + " XXX% free"

var (
	green  = lipgloss.Color("2")
	cyan   = lipgloss.Color("6")
	yellow = lipgloss.Color("3")
	red    = lipgloss.Color("1")
	white  = lipgloss.Color("15")
	gray   = lipgloss.Color("237")

	headerStyle  = lipgloss.NewStyle().Bold(true).Foreground(white)
	sectionStyle = lipgloss.NewStyle().Bold(true).Foreground(white).Underline(true)
	dimStyle     = lipgloss.NewStyle().Faint(true)
	labelStyle   = lipgloss.NewStyle().Bold(true).Foreground(white)
)

func barFilledStyle(c lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Background(c).Foreground(white)
}

var barEmptyStyle = lipgloss.NewStyle().Background(gray)

var (
	modelBarFilledStyle = lipgloss.NewStyle().Foreground(cyan)
	modelBarEmptyStyle  = lipgloss.NewStyle().Foreground(gray)
)

func pctStyle(c lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c)
}

func usageColor(pct float64) lipgloss.Color {
	switch {
	case pct >= 90:
		return red
	case pct >= 70:
		return yellow
	default:
		return green
	}
}
