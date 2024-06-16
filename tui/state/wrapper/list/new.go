package list

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
)

func New(list list.Model, other help.KeyMap) *State {
	return &State{
		list:   list,
		keyMap: newKeyMap(list.KeyMap, other),
	}
}
