package path

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		Copy: util.Bind("copy", "c", "enter"),
		quit: util.Bind("quit", "q", "ctrl+c"),
	}
}

// keyMap implements help.KeyMap.
type keyMap struct {
	Copy,
	quit key.Binding
}

// ShortHelp implements help.KeyMap.
func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Copy,
		k.quit,
	}
}

// FullHelp implements help.KeyMap.
func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
