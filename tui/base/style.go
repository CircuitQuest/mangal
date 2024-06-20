package base

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

const (
	Ellipsis         = "…"
	StatusSeparator  = " • "
	HelpKeySeparator = StatusSeparator
)

var SpinnerType = spinner.Dot

type styles struct {
	title,
	status,
	notification,
	subtitle,
	header,
	state,
	footer lipgloss.Style

	loading loadingStyle
	help    help.Styles
}

type loadingStyle struct {
	spinner,
	message lipgloss.Style
}

func defaultStyles() styles {
	// Flipped Accent
	tempAccent := style.Bold.Accent.Background(color.Background)

	helpKey := style.Bold.Warning
	helpSep := style.Normal.Base
	helpDesc := style.Normal.Secondary

	return styles{
		title:        style.FlipGrounds(tempAccent).Padding(0, 1).Margin(0, 0, 0, 1),
		status:       style.Normal.Base.Padding(0, 0, 0, 1),
		notification: style.Normal.Warning.Padding(0, 0, 0, 1),
		subtitle:     style.Normal.Secondary.Padding(0, 0, 0, 1),
		header:       style.Normal.Base.Padding(0, 0, 1, 1),
		state:        style.Normal.Base.Padding(0, 1),
		footer:       style.Normal.Base.Padding(0, 1),
		loading: loadingStyle{
			spinner: style.Bold.Accent,
			message: style.Normal.Secondary.Padding(0, 0, 0, 1),
		},
		help: help.Styles{
			Ellipsis:       style.Normal.Base,
			ShortKey:       helpKey,
			ShortDesc:      helpDesc,
			ShortSeparator: helpSep,
			FullKey:        helpKey,
			FullDesc:       helpDesc,
			FullSeparator:  helpSep,
		},
	}
}
