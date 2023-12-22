package listwrapper

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/util"
)

func New(list list.Model) *State {
	return &State{
		list: list,
		keyMap: KeyMap{
			reverse: util.Bind("reverse", "R"),
			list:    &list.KeyMap,
		},
	}
}
