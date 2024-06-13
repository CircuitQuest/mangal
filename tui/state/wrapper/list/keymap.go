package list

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

var _ help.KeyMap = (*keyMap)(nil)

// keyMap implements help.keyMap.
type keyMap struct {
	reverse key.Binding

	list *list.KeyMap
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.list.Filter,
		k.reverse,
	}
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
