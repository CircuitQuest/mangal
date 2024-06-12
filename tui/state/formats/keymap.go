package formats

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

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
