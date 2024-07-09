package mangas

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		confirm:       util.Bind("confirm", "enter"),
		search:        util.Bind("search", "s"),
		cancelSearch:  util.Bind("cancel search", "esc"),
		confirmSearch: util.Bind("confirm search", "enter"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	confirm,
	search,
	cancelSearch,
	confirmSearch key.Binding
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.confirm,
		k.search,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
