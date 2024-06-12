package formats

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
)

func New() *State {
	listWrapper := listwrapper.New(util.NewList(
		2,
		"manga", "mangas",
		libmangal.FormatValues(),
		func(format libmangal.Format) list.DefaultItem {
			return Item{format: format}
		},
	))

	return &State{
		list: listWrapper,
		keyMap: keyMap{
			setRead:     util.Bind("set for reading", "r"),
			setDownload: util.Bind("set for downloading", "d"),
			setAll:      util.Bind("set for all", "enter"),
			list:        listWrapper.KeyMap(),
		},
	}
}
