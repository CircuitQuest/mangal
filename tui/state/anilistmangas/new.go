package anilistmangas

import (
	"github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
)

func New(anilist *lmanilist.Anilist, chapters []lmanilist.Manga, onResponse OnResponseFunc) *State {
	listWrapper := listwrapper.New(util.NewList(
		2,
		"manga", "mangas",
		chapters,
		func(manga lmanilist.Manga) list.DefaultItem {
			return Item{Manga: &manga}
		},
	))

	return &State{
		anilist:    anilist,
		list:       listWrapper,
		onResponse: onResponse,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
			Search:  util.Bind("search", "s"),
		},
	}
}
