package providers

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		info:    util.Bind("info", "i"),
		confirm: util.Bind("confirm", "enter"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	info,
	confirm key.Binding
}

// ShortHelp implements help.KeyMap.
func (p keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		p.confirm,
		p.info,
	}
}

// FullHelp implements help.KeyMap.
func (p keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		p.ShortHelp(),
	}
}
