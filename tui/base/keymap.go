package base

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var (
	_ help.KeyMap = (*keyMap)(nil)
	_ help.KeyMap = (*combinedKeyMap)(nil)
	_ help.KeyMap = (*NoKeyMap)(nil)
)

func newKeyMap() *keyMap {
	return &keyMap{
		back: util.Bind("back", "esc"),
		quit: util.Bind("quit", "ctrl+c"),
		help: util.Bind("help", "?"),
		log:  util.Bind("log", "ctrl+l"),
	}
}

// keyMap implements help.KeyMap.
type keyMap struct {
	back,
	quit,
	help,
	log key.Binding
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.back,
		k.help,
	}
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{
			k.quit,
			k.log,
		},
	}
}

// with creates a combined keymap
func (k keyMap) with(other help.KeyMap) combinedKeyMap {
	return combinedKeyMap{
		k:     k,
		other: other,
	}
}

// combinedKeyMap implements help.KeyMap.
type combinedKeyMap struct {
	k     help.KeyMap
	other help.KeyMap
}

// ShortHelp implements help.KeyMap.
func (c combinedKeyMap) ShortHelp() []key.Binding {
	return append(c.k.ShortHelp(), c.other.ShortHelp()...)
}

// FullHelp implements help.KeyMap.
func (c combinedKeyMap) FullHelp() [][]key.Binding {
	return append(c.k.FullHelp(), c.other.FullHelp()...)
}

// NoKeyMap implements help.keyMap.
type NoKeyMap struct{}

// ShortHelp implements help.keyMap.
func (k NoKeyMap) ShortHelp() []key.Binding {
	return nil
}

// FullHelp implements help.keyMap.
func (k NoKeyMap) FullHelp() [][]key.Binding {
	return nil
}
