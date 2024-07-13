package anilist

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(anilist *lmanilist.Anilist, manga mangadata.Manga) *state {
	_keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"anilist manga", "anilist mangas",
		nil,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return &item{manga: manga}
		},
		&_keyMap)

	title := manga.Info().AnilistSearch
	if title == "" {
		title = manga.Info().Title
	}
	s := &state{
		anilist: anilist,
		search:  search.New("Search anilist manga...", title),
		manga:   manga,
		list:    listWrapper,
		keyMap:  &_keyMap,
	}
	s.updateKeybinds()
	return s
}
