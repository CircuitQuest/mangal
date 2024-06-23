package providers

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/zyedidia/generic/set"
)

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	loader      libmangal.ProviderLoader
	loadedItems *set.Set[*item]
	loadTime    time.Time
	extraInfo   *bool
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	return i.loader.String()
}

// Title implements list.DefaultItem.
func (i *item) Title() string {
	var title strings.Builder
	title.WriteString(i.FilterValue())

	if i.isLoaded() {
		title.WriteString(" ")
		title.WriteString(icon.Check.String())
		if *i.extraInfo {
			timeAgo := fmt.Sprintf(" loaded %s ago", time.Since(i.loadTime).Truncate(time.Second).String())
			title.WriteString(style.Normal.Secondary.Render(timeAgo))
		}
	}

	return title.String()
}

// Description implements list.DefaultItem.
func (i *item) Description() string {
	info := i.loader.Info()
	if *i.extraInfo {
		moreInfo := fmt.Sprintf("%s v%s", info.ID, info.Version)
		return lipgloss.JoinVertical(lipgloss.Left, moreInfo, info.Website)
	}

	return info.Website
}

func (i *item) isLoaded() bool {
	return i.loadedItems.Has(i)
}

func (i *item) markLoaded() {
	if !i.isLoaded() {
		i.loadTime = time.Now()
		i.loadedItems.Put(i)
	}
}

func (i *item) markClosed() {
	if i.isLoaded() {
		i.loadedItems.Remove(i)
	}
}
