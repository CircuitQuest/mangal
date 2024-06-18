package base

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type styles struct {
	title,
	status,
	notification,
	subtitle,
	header,
	state,
	footer,
	spinner,
	loading,
	helpKey,
	helpSep lipgloss.Style
}

func defaultStyles() styles {
	// Flipped Accent
	tempAccent := style.Bold.Accent.Background(color.Background)
	return styles{
		title:        style.FlipGrounds(tempAccent).Padding(0, 1),
		status:       style.Normal.Base.Padding(0, 0, 0, 1),
		notification: style.Normal.Warning.Padding(0, 0, 0, 1),
		subtitle:     style.Normal.Secondary.Padding(1, 0, 0, 0),
		header:       style.Normal.Base.Padding(0, 0, 1, 2),
		state:        style.Normal.Base.Padding(0, 1),
		footer:       style.Normal.Base.Padding(0, 1),
		spinner:      style.Bold.Accent,
		loading:      style.Normal.Secondary.Padding(0, 0, 0, 1),
		helpKey:      style.Bold.Warning,
		helpSep:      style.Normal.Base,
	}
}
