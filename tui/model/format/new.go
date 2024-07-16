package format

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/help"
	"github.com/luevano/mangal/tui/model/list"
)

func New(accentColor lipgloss.Color) *Model {
	listWrapper := list.New(
		2, "manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) _list.DefaultItem {
			return &item{format: format}
		},
	)

	t := lipgloss.NewStyle().
		Background(accentColor).
		Padding(0, 1).
		Margin(0, 1, 1, 1).
		Render("Formats")
	h := help.New()
	h.ShowAll = true
	_keyMap := newKeyMap()
	return &Model{
		list:  listWrapper,
		title: t,
		help:  h.View(_keyMap),
		size: base.Size{ // basically perfect size, will make tweakable later
			Width:  24,
			Height: 18,
		},
		keyMap: _keyMap,
	}
}
