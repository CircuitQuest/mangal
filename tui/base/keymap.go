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
		quit: util.Bind("quit", "ctrl+c"),
		back: util.Bind("back", "esc"),
		home: util.Bind("home", "H"),
		help: util.Bind("help", "?"),
		log:  util.Bind("log", "ctrl+l"),
	}
}

// keyMap implements help.KeyMap.
type keyMap struct {
	quit,
	back,
	home,
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
			k.home,
			k.quit,
			k.log,
		},
	}
}

// CombinedKeyMap is a convenience function to create a combined KeyMap that
// implements help.KeyMap, useful for concatenating them.
//
// The resulting keybind is not ment to be be modified in-place.
func CombinedKeyMap(first, second help.KeyMap) combinedKeyMap {
	return combinedKeyMap{
		first:  first,
		second: second,
	}
}

// combinedKeyMap implements help.KeyMap.
type combinedKeyMap struct {
	first  help.KeyMap
	second help.KeyMap
}

// ShortHelp implements help.KeyMap.
func (c combinedKeyMap) ShortHelp() []key.Binding {
	return append(c.first.ShortHelp(), c.second.ShortHelp()...)
}

// FullHelp implements help.KeyMap.
func (c combinedKeyMap) FullHelp() [][]key.Binding {
	return append(c.first.FullHelp(), c.second.FullHelp()...)
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
