package errorstate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

// keyMap implements help.keyMap.
type keyMap struct {
	quit,
	copyError key.Binding
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.quit,
		k.copyError,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
