package anilistmangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(anilist *lmanilist.Anilist, mangas []lmanilist.Manga, onResponse onResponseFunc) *state {
	_keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"anilist manga", "anilist mangas",
		mangas,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return &item{manga: manga}
		},
		_keyMap)

	s := &state{
		anilist:    anilist,
		search:     search.New("Search anilist manga...", ""),
		list:       listWrapper,
		onResponse: onResponse,
		keyMap:     _keyMap,
	}

	return s
}
