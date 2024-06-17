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
	_ list.Item        = (*Item)(nil)
	_ list.DefaultItem = (*Item)(nil)
)

// Item implements list.Item.
type Item struct {
	libmangal.ProviderLoader
	loadedItems *set.Set[*Item]
	loadTime    time.Time
	extraInfo   *bool
}

// FilterValue implements list.Item.
func (i *Item) FilterValue() string {
	return i.String()
}

// Title implements list.DefaultItem.
func (i *Item) Title() string {
	var title strings.Builder
	title.WriteString(i.FilterValue())

	if i.IsLoaded() {
		title.WriteString(" ")
		title.WriteString(style.Bold.Success.Render(icon.Check.String()))
		if *i.extraInfo {
			timeAgo := fmt.Sprintf(" loaded %s ago", time.Since(i.loadTime).Truncate(time.Second).String())
			title.WriteString(style.Normal.Secondary.Render(timeAgo))
		}
	}

	return title.String()
}

// Description implements list.DefaultItem.
func (i *Item) Description() string {
	info := i.Info()
	if *i.extraInfo {
		return leftExtraInfo(info)
	}

	return info.Website
}

func (i *Item) IsLoaded() bool {
	return i.loadedItems.Has(i)
}

func (i *Item) MarkLoaded() {
	if !i.IsLoaded() {
		i.loadTime = time.Now()
		i.loadedItems.Put(i)
	}
}

func (i *Item) MarkClosed() {
	if i.IsLoaded() {
		i.loadedItems.Remove(i)
	}
}

func leftExtraInfo(info libmangal.ProviderInfo) string {
	moreInfo := fmt.Sprintf("%s v%s", info.ID, info.Version)
	return lipgloss.JoinVertical(lipgloss.Left, moreInfo, info.Website)
}
