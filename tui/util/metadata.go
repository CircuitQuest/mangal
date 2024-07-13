package util

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/color"
)

type MetaStyle struct {
	Color  lipgloss.Color
	Prefix string
	Code   string
}

func MetaIDStyle(id metadata.ID) MetaStyle {
	c := color.Provider
	p := "Provider"
	v := id.Code
	switch id.Source {
	case metadata.IDSourceAnilist:
		c = color.Anilist
		p = "Anilist"
		v = metadata.IDCodeAnilist
	case metadata.IDSourceMyAnimeList:
		c = color.MyAnimeList
		p = "MyAnimeList"
		v = metadata.IDCodeMyAnimeList
	case metadata.IDSourceKitsu:
		c = color.Kitsu
		p = "Kitsu"
		v = metadata.IDCodeKitsu
	case metadata.IDSourceMangaUpdates:
		c = color.MangaUpdates
		p = "MangaUpdates"
		v = metadata.IDCodeMangaUpdates
	case metadata.IDSourceAnimePlanet:
		c = color.AnimePlanet
		p = "AnimePlanet"
		v = metadata.IDCodeAnimePlanet
	}
	return MetaStyle{c, p, v}
}
