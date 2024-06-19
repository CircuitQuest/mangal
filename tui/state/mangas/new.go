package mangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(client *libmangal.Client, query string, mangas []mangadata.Manga) *State {
	keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"manga", "mangas",
		mangas,
		func(manga mangadata.Manga) _list.DefaultItem {
			return &Item{manga}
		},
		keyMap)

	return &State{
		list:   listWrapper,
		mangas: mangas,
		client: client,
		query:  query,
		keyMap: keyMap,
	}
}
