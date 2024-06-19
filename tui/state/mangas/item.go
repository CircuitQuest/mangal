package mangas

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal/mangadata"
)

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	manga mangadata.Manga
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	return i.manga.String()
}

// Title implements list.DefaultItem.
func (i *item) Title() string {
	return i.FilterValue()
}

// Description implements list.DefaultItem.
func (i *item) Description() string {
	return i.manga.Info().URL
}
