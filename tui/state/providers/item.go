package providers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
)

var (
	_ list.Item        = (*Item)(nil)
	_ list.DefaultItem = (*Item)(nil)
)

// Item implements list.Item.
type Item struct {
	libmangal.ProviderLoader
}

// FilterValue implements list.Item.
func (i Item) FilterValue() string {
	return i.String()
}

// Title implements list.DefaultItem.
func (i Item) Title() string {
	return i.FilterValue()
}

// Description implements list.DefaultItem.
func (i Item) Description() string {
	info := i.Info()
	return fmt.Sprintf(
		"%s v%s\n%s",
		info.ID,
		info.Version,
		info.Website,
	)
}
