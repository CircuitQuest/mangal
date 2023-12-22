package mangas

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/state/listwrapper"
)

var _ help.KeyMap = (*KeyMap)(nil)

type KeyMap struct {
	Confirm key.Binding
	list    listwrapper.KeyMap
}

func (k KeyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.Confirm,
	)
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
