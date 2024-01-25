package base

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/ui/color"
)

type Styles struct {
	Title,
	TitleBar,
	Subtitle,
	HelpBar,
	Selection lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Title: lipgloss.
			NewStyle().
			Bold(true).
			Background(color.Accent).
			Foreground(color.Background).
			Padding(0, 1),
		TitleBar: lipgloss.
			NewStyle().
			Padding(0, 0, 1, 2),
		Subtitle: lipgloss.
			NewStyle().
			Foreground(color.Secondary),
		HelpBar: lipgloss.
			NewStyle().
			Padding(0, 1),
	}
}
