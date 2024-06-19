package volumes

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal/mangadata"
)

var (
	_ list.Item        = (*Item)(nil)
	_ list.DefaultItem = (*Item)(nil)
)

// Item implements list.Item.
type Item struct {
	volume mangadata.Volume
}

// FilterValue implements list.Item.
func (i *Item) FilterValue() string {
	return fmt.Sprintf("Volume %s", i.volume)
}

// Title implements list.DefaultItem.
func (i *Item) Title() string {
	return i.FilterValue()
}

// Description implements list.DefaultItem.
func (i *Item) Description() string {
	return ""
}
