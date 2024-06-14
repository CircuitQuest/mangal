package mangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New(client *libmangal.Client, query string, mangas []mangadata.Manga) *State {
	listWrapper := list.New(util.NewList(
		2,
		"manga", "mangas",
		mangas,
		func(manga mangadata.Manga) _list.DefaultItem {
			return Item{manga}
		},
	))

	return &State{
		list:   listWrapper,
		mangas: mangas,
		client: client,
		query:  query,
		keyMap: newKeyMap(listWrapper.KeyMap()),
	}
}
