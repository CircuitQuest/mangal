package list

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(listKeyMap list.KeyMap, other help.KeyMap) keyMap {
	return keyMap{
		reverse: util.Bind("reverse", "R"),
		list:    listKeyMap,
		other:   other,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	reverse key.Binding
	list    list.KeyMap

	other help.KeyMap
}

func (k keyMap) shortHelp() []key.Binding {
	return []key.Binding{
		k.list.Filter,
		k.reverse,
		k.list.CursorUp,
		k.list.CursorDown,
	}
}

// ShortHelp implements help.KeyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return append(k.other.ShortHelp(),
		k.shortHelp()...,
	)
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return append(k.other.FullHelp(),
		[][]key.Binding{
			k.shortHelp(),
			{
				k.list.NextPage,
				k.list.PrevPage,
				k.list.GoToStart,
				k.list.GoToEnd,
			},
			{
				k.list.ClearFilter,
				k.list.CancelWhileFiltering,
				k.list.AcceptWhileFiltering,
			},
		}...,
	)
}
