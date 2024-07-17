package mangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/tui/model/search"
)

func New(client *libmangal.Client) *state {
	info := false
	fullInfo := false
	listWrapper := list.New(
		2,
		"manga", "mangas",
		nil,
		func(manga mangadata.Manga) _list.DefaultItem {
			return newItem(manga, &info, &fullInfo)
		},
	)

	s := &state{
		list:          listWrapper,
		search:        search.New("Search manga...", "", 64, 5),
		client:        client,
		extraInfo:     &info,
		fullExtraInfo: &fullInfo,
		keyMap:        newKeyMap(),
	}
	s.updateKeybinds()
	return s
}
