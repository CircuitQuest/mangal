package anilistmangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(anilist *lmanilist.Anilist, chapters []lmanilist.Manga, onResponse OnResponseFunc) *State {
	keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"manga", "mangas",
		chapters,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return &Item{Manga: &manga}
		},
		keyMap)

	return &State{
		anilist:    anilist,
		list:       listWrapper,
		onResponse: onResponse,
		keyMap:     keyMap,
	}
}
