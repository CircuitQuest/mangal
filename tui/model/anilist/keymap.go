package anilist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		login:   util.Bind("login", "i"),
		open:    util.Bind("open auth url", "ctrl+o"),
		logout:  util.Bind("logout", "o"),
		up:      util.BindNamedKey("↑/k", "up", "k", "up"),
		down:    util.BindNamedKey("↓/j", "down", "j", "down"),
		selekt:  util.Bind("select", "enter"),
		confirm: util.Bind("confirm", "enter"),
		cancel:  util.Bind("cancel", "esc"),
		clear:   util.Bind("clear field", "c"),
		quit:    util.Bind("quit", "q", "ctrl+c"),
	}
}

// keyMap implements help.KeyMap.
type keyMap struct {
	login,
	open,
	logout,
	up,
	down,
	selekt, // select input
	confirm, // confirm input content
	cancel, // cancel input content
	clear, // clear input content
	quit key.Binding
}

// ShortHelp implements help.KeyMap.
func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.login,
		k.open,
		k.logout,
		k.selekt,
		k.clear,
		k.up,
		k.down,
		k.quit,
	}
}

// FullHelp implements help.KeyMap.
func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.login,
			k.open,
			k.logout,
		},
		{
			k.selekt,
			k.clear,
		},
		{
			k.up,
			k.down,
		},
		{
			k.quit,
		},
	}
}
