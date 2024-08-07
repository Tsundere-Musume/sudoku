package main

import "github.com/charmbracelet/lipgloss"

const ()

var (
	//Colors
	active = lipgloss.Color("#f6c177")
	// inactive     = lipgloss.Color("#c4a7e7")
	invalid      = lipgloss.Color("#eb6f92")
	valid        = lipgloss.Color("#a6e3a1")
	borderColor  = lipgloss.Color("#31748f")
	TEXT         = lipgloss.Color("#313244")
	NOT_EDITABLE = lipgloss.Color("#ebbcba")

	base = lipgloss.NewStyle().Padding(0, 1).Foreground(TEXT).BorderForeground(borderColor)

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       " ",
		TopLeft:     "╭",
		TopRight:    "┘",
		BottomLeft:  "╰",
		BottomRight: "┐",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "┤",
		BottomLeft:  "╰",
		BottomRight: "┤",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		Padding(0, 1).
		BorderForeground(borderColor)

	activeTab = tab.Border(activeTabBorder, true).Foreground(active)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderBottom(false)

	borderStyle = lipgloss.NewStyle().Foreground(borderColor).BorderForeground(borderColor).Padding(0, 1)
	boardBorder = borderStyle.UnsetPadding().BorderStyle(lipgloss.ThickBorder())
)
