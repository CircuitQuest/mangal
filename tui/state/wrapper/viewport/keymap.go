package viewport

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/luevano/mangal/tui/util"
)

// TODO: add custom keybinds

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(vieportKeyMap *viewport.KeyMap) keyMap {
	return keyMap{
		viewport: vieportKeyMap,
		copy:     util.Bind("copy content", "c"),
		goTop:    util.BindNamedKey("g/home", "go to start", "g", "home"),
		goBottom: util.BindNamedKey("G/end", "go to end", "G", "end"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	viewport *viewport.KeyMap

	copy,
	goTop,
	goBottom key.Binding
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.copy,
		k.viewport.Up,
		k.viewport.Down,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{
			k.goTop,
			k.goBottom,
		},
		{
			k.viewport.HalfPageUp,
			k.viewport.HalfPageDown,
		},
		{
			k.viewport.PageUp,
			k.viewport.PageDown,
		},
	}
}
