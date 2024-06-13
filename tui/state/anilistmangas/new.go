package anilistmangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New(anilist *lmanilist.Anilist, chapters []lmanilist.Manga, onResponse OnResponseFunc) *State {
	listWrapper := list.New(util.NewList(
		2,
		"manga", "mangas",
		chapters,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return Item{Manga: &manga}
		},
	))

	return &State{
		anilist:    anilist,
		list:       listWrapper,
		onResponse: onResponse,
		keyMap: keyMap{
			confirm: util.Bind("confirm", "enter"),
			search:  util.Bind("search", "s"),
		},
	}
}
