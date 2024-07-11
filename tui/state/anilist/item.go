package anilist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
)

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	manga lmanilist.Manga
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	title := i.manga.String()
	if title != "" {
		return title
	}
	return "Untitled"
}

// Item implements list.Item.
func (i *item) Title() string {
	return i.FilterValue()
}

// Description implements list.Item.
func (i *item) Description() string {
	return fmt.Sprint("https://anilist.co/manga/", i.manga.ID)
}
