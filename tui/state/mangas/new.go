package mangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(client *libmangal.Client) *state {
	_keyMap := newKeyMap()
	info := false
	fullInfo := false
	listWrapper := list.New(
		2,
		"manga", "mangas",
		nil,
		func(manga mangadata.Manga) _list.DefaultItem {
			return newItem(manga, &info, &fullInfo)
		},
		&_keyMap)

	s := &state{
		list:          listWrapper,
		search:        search.New("Search manga...", ""),
		client:        client,
		extraInfo:     &info,
		fullExtraInfo: &fullInfo,
		keyMap:        &_keyMap,
	}
	s.updateKeybinds()

	return s
}
