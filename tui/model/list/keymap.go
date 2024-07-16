package list

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*KeyMap)(nil)

func newKeyMap(listKeyMap *list.KeyMap, other help.KeyMap) KeyMap {
	return KeyMap{
		List:    listKeyMap,
		other:   other,
		Reverse: util.Bind("reverse", "R"),
	}
}

// KeyMap implements help.KeyMap.
type KeyMap struct {
	List  *list.KeyMap
	other help.KeyMap

	Reverse key.Binding
}

func (k KeyMap) shortHelp() []key.Binding {
	return []key.Binding{
		k.List.Filter,
		k.Reverse,
		k.List.CursorUp,
		k.List.CursorDown,
	}
}

// ShortHelp implements help.KeyMap.
func (k KeyMap) ShortHelp() []key.Binding {
	return append(k.other.ShortHelp(),
		k.shortHelp()...,
	)
}

// FullHelp implements help.KeyMap.
func (k KeyMap) FullHelp() [][]key.Binding {
	return append(k.other.FullHelp(),
		[][]key.Binding{
			k.shortHelp(),
			{
				k.List.NextPage,
				k.List.PrevPage,
				k.List.GoToStart,
				k.List.GoToEnd,
			},
			{
				k.List.ClearFilter,
				k.List.CancelWhileFiltering,
				k.List.AcceptWhileFiltering,
			},
		}...,
	)
}
