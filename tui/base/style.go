package base

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type styles struct {
	title,
	notification,
	subtitle,
	header,
	footer lipgloss.Style
}

func defaultStyles() styles {
	// Flipped Accent
	tempAccent := style.Bold.Accent.Background(color.Background)
	return styles{
		title:        style.FlipGrounds(tempAccent).Padding(0, 1),
		notification: style.Normal.Warning.Padding(0, 0, 0, 1),
		subtitle:     style.Normal.Secondary,
		header:       style.Normal.Base.Padding(0, 0, 1, 2),
		footer:       style.Normal.Base.Padding(0, 1),
	}
}
