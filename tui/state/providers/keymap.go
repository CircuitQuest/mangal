package providers

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

// keyMap implements help.keyMap.
type keyMap struct {
	info,
	confirm key.Binding

	list help.KeyMap
}

// ShortHelp implements help.KeyMap.
func (p keyMap) ShortHelp() []key.Binding {
	return append(
		p.list.ShortHelp(),
		p.confirm,
		p.info,
	)
}

// FullHelp implements help.KeyMap.
func (p keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		p.ShortHelp(),
	}
}
