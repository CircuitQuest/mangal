package formats

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		setRead:     util.Bind("set for reading", "r"),
		setDownload: util.Bind("set for downloading", "d"),
		setBoth:     util.Bind("set for both", "enter"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	setRead,
	setDownload,
	setBoth key.Binding
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.setRead,
		k.setDownload,
	}
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		append(k.ShortHelp(), k.setBoth),
	}
}
