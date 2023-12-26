package chapters

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
	"github.com/luevano/libmangal"
	"github.com/zyedidia/generic/set"
)

func New(client *libmangal.Client, volume libmangal.Volume, chapters []libmangal.Chapter) *State {
	selectedSet := set.NewMapset[*Item]()
	listWrapper := listwrapper.New(util.NewList(
		2,
		"chapter", "chapters",
		chapters,
		func(chapter libmangal.Chapter) list.DefaultItem {
			return &Item{
				chapter:       chapter,
				selectedItems: &selectedSet,
				client:        client,
			}
		},
	))

	return &State{
		client:   client,
		volume:   volume,
		selected: selectedSet,
		list:     listWrapper,
		keyMap: KeyMap{
			UnselectAll:  util.Bind("unselect all", "backspace"),
			SelectAll:    util.Bind("select all", "a"),
			Toggle:       util.Bind("toggle", " "),
			Read:         util.Bind("read", "r"),
			OpenURL:      util.Bind("open url", "o"),
			Anilist:      util.Bind("anilist", "A"),
			Download:     util.Bind("download", "d"),
			Confirm:      util.Bind("confirm", "enter"),
			ChangeFormat: util.Bind("change format", "f"),
			list:         listWrapper.GetKeyMap(),
		},
	}
}
