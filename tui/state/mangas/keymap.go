package mangas

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(listKeyMap help.KeyMap) keyMap {
	return keyMap{
		confirm: util.Bind("confirm", "enter"),
		list:    listKeyMap,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	confirm key.Binding

	list help.KeyMap
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.confirm,
	)
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
