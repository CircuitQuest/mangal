package anilist

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/tui/model/search"
)

func New(anilist *lmanilist.Anilist, manga mangadata.Manga) *state {
	listWrapper := list.New(
		2,
		"anilist manga", "anilist mangas",
		nil,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return &item{manga: manga}
		},
	)

	title := manga.Info().AnilistSearch
	if title == "" {
		title = manga.Info().Title
	}
	s := &state{
		anilist: anilist,
		search:  search.New("Search anilist manga...", title, 64, 5),
		manga:   manga,
		list:    listWrapper,
		keyMap:  newKeyMap(),
	}
	s.updateKeybinds()
	return s
}
