package chapsdownloaded

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(state *State) keyMap {
	return keyMap{
		open:  util.Bind("open directory", "o"),
		quit:  util.Bind("quit", "q"),
		retry: util.Bind("retry", "r"),
		state: state,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	quit,
	open,
	retry key.Binding

	state *State
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	bindings := []key.Binding{
		k.quit,
		k.open,
	}

	if len(k.state.options.Failed) > 0 {
		bindings = append(bindings, k.retry)
	}

	return bindings
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
	}
}
