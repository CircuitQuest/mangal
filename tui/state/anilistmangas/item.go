package anilistmangas

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
)

var (
	_ list.Item        = (*Item)(nil)
	_ list.DefaultItem = (*Item)(nil)
)

// Item implements list.Item.
type Item struct {
	Manga *lmanilist.Manga
}

// FilterValue implements list.Item.
func (i Item) FilterValue() string {
	for _, title := range []string{
		i.Manga.Title.English,
		i.Manga.Title.Romaji,
		i.Manga.Title.Native,
	} {
		if title != "" {
			return title
		}
	}

	return "Untitled"
}

// Item implements list.Item.
func (i Item) Title() string {
	return i.FilterValue()
}

// Description implements list.Item.
func (i Item) Description() string {
	return fmt.Sprint("https://anilist.co/manga/", i.Manga.ID)
}
