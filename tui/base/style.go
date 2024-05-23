package base

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type Styles struct {
	Title,
	TitleBar,
	Subtitle,
	HelpBar,
	Selection lipgloss.Style
}

func DefaultStyles() Styles {
	// Flipped Accent
	tempAccent := style.Bold.Accent.Background(color.Background)
	return Styles{
		Title:    style.FlipGrounds(tempAccent).Padding(0, 1),
		TitleBar: style.Normal.Base.Padding(0, 0, 1, 2),
		Subtitle: style.Normal.Secondary,
		HelpBar:  style.Normal.Base.Padding(0, 1),
	}
}
