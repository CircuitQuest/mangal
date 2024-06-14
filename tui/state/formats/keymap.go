package formats

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(listKeyMap help.KeyMap) keyMap {
	return keyMap{
		setRead:     util.Bind("set for reading", "r"),
		setDownload: util.Bind("set for downloading", "d"),
		setAll:      util.Bind("set for all", "enter"),
		list:        listKeyMap,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	setRead,
	setDownload,
	setAll key.Binding

	list help.KeyMap
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.setRead,
		k.setDownload,
	)
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
