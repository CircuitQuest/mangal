package formats

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New() *State {
	listWrapper := list.New(util.NewList(
		2,
		"manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) _list.DefaultItem {
			return Item{format: format}
		},
	))

	return &State{
		list:   listWrapper,
		keyMap: newKeyMap(listWrapper.KeyMap()),
	}
}
