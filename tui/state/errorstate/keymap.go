package errorstate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		copyError: util.Bind("copy error", "c"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	copyError key.Binding
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.copyError,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
