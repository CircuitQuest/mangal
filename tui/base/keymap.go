package base

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*NoKeyMap)(nil)

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

type keyMap struct {
	back,
	quit,
	help,
	log key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		back: util.Bind("back", "esc"),
		quit: util.Bind("quit", "ctrl+c"),
		help: util.Bind("help", "?"),
		log:  util.Bind("log", "ctrl+l"),
	}
}
