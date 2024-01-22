package model

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

type keyMap struct {
	Back,
	Quit,
	Help,
	Log key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		Back: util.Bind("back", "esc"),
		Quit: util.Bind("quit", "ctrl+c"),
		Help: util.Bind("help", "?"),
		Log:  util.Bind("log", "ctrl+l"),
	}
}
