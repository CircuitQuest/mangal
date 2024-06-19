package formats

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New() *state {
	keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) _list.DefaultItem {
			return &item{format: format}
		},
		keyMap)

	return &state{
		list:   listWrapper,
		keyMap: keyMap,
	}
}
