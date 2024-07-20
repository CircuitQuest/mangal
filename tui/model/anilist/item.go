package anilist

import "github.com/charmbracelet/bubbles/list"

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	user string
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	return i.user
}

// Item implements list.DefaultItem.
func (i *item) Title() string {
	return i.FilterValue()
}

// Description implements list.DefaultItem.
func (i *item) Description() string {
	return ""
}
