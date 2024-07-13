package metadata

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	_style "github.com/luevano/mangal/theme/style"
)

type style struct {
	Color  lipgloss.Color
	Prefix string
}

type styles struct {
	base lipgloss.Style

	provider,
	anilist,
	myAnimeList,
	kitsu,
	mangaUpdates,
	animePlanet style
}

func defaultStyles() styles {
	return styles{
		base:         _style.Normal.Base.Padding(0, 1).Foreground(color.Bright),
		provider:     style{color.Provider, "Provider"},
		anilist:      style{color.Anilist, "Anilist"},
		myAnimeList:  style{color.MyAnimeList, "MyAnimeList"},
		kitsu:        style{color.Kitsu, "Kitsu"},
		mangaUpdates: style{color.MangaUpdates, "MangaUpdates"},
		animePlanet:  style{color.AnimePlanet, "AnimePlanet"},
	}
}
