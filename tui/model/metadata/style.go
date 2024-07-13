package metadata

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/metadata"
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
		anilist:      style{color.Anilist, "Anilist(" + metadata.IDCodeAnilist + ")"},
		myAnimeList:  style{color.MyAnimeList, "MyAnimeList(" + metadata.IDCodeMyAnimeList + ")"},
		kitsu:        style{color.Kitsu, "Kitsu(" + metadata.IDCodeKitsu + ")"},
		mangaUpdates: style{color.MangaUpdates, "MangaUpdates(" + metadata.IDCodeMangaUpdates + ")"},
		animePlanet:  style{color.AnimePlanet, "AnimePlanet(" + metadata.IDCodeAnimePlanet + ")"},
	}
}
