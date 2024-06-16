package viewport

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

// TODO: add custom keybinds

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(vieportKeyMap viewport.KeyMap) keyMap {
	return keyMap{
		viewport: vieportKeyMap,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	viewport viewport.KeyMap
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.viewport.Up,
		k.viewport.Down,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{
			k.viewport.HalfPageUp,
			k.viewport.HalfPageDown,
		},
		{
			k.viewport.PageUp,
			k.viewport.PageDown,
		},
	}
}
