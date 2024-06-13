package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/util"
)

func New(list list.Model) *State {
	return &State{
		list: list,
		keyMap: keyMap{
			reverse: util.Bind("reverse", "R"),
			list:    &list.KeyMap,
		},
	}
}
