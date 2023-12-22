package anilistmangas

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
	"github.com/mangalorg/libmangal"
)

func New(anilist *libmangal.Anilist, chapters []libmangal.AnilistManga, onResponse OnResponseFunc) *State {
	listWrapper := listwrapper.New(util.NewList(
		2,
		"manga", "mangas",
		chapters,
		func(manga libmangal.AnilistManga) list.DefaultItem {
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
