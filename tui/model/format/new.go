package format

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/help"
	"github.com/luevano/mangal/tui/model/list"
)

func New(accentColor lipgloss.Color) *Model {
	list := list.New(
		2, 0,
		"manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) _list.DefaultItem {
			return &item{format: format}
		},
	)
	list.SetAccentColor(accentColor)

	t := lipgloss.NewStyle().
		Background(accentColor).
		Padding(0, 1).
		Margin(0, 1, 1, 1).
		Render("Formats")
	// needs to be rendered here, as rendering aat item.go
	// (on its definition) will load the incorrect icon type
	sep = style.Bold.Warning.Padding(0, 1).Render(icon.Separator.Raw())

	h := help.New()
	h.ShowAll = true
	_keyMap := newKeyMap()
	return &Model{
		list:  list,
		title: t,
		help:  h.View(_keyMap),
		size: base.Size{ // basically perfect size, will make tweakable later
			Width:  24,
			Height: 12, // only space required for current formats
		},
		keyMap: _keyMap,
	}
}
