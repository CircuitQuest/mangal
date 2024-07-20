package list

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/icon"
)

func New[T any](
	itemHeight int,
	itemSpacing int,
	singular, plural string,
	items []T,
	transform func(T) list.DefaultItem,
) *Model {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = transform(item)
	}

	border := lipgloss.ThickBorder()
	delegate := list.NewDefaultDelegate()

	// TODO: possibly use the current "window" (where the list is being displayed) accent color,
	// instead of always hardcoding color.Accent
	//
	// Styles don't use mangal/theme/style, as they are more specialized with paddings and whatnot
	styles := delegate.Styles
	styles.NormalTitle = styles.NormalTitle.Bold(true)
	styles.SelectedTitle = styles.SelectedTitle.Bold(true).
		Foreground(color.Accent).
		Border(border, false, false, false, true).
		BorderForeground(color.Accent)
	styles.SelectedDesc = styles.SelectedDesc.
		Foreground(delegate.Styles.NormalDesc.GetForeground()).
		Border(border, false, false, false, true).
		BorderForeground(color.Accent)
	delegate.Styles = styles

	if itemHeight == 1 {
		delegate.ShowDescription = false
	}
	delegate.SetHeight(itemHeight)
	delegate.SetSpacing(itemSpacing)

	// If only a few items, just set the necessary hight, else the max
	perItemHeight := itemHeight + itemSpacing
	itemsHeight := perItemHeight * len(items)

	l := list.New(listItems, delegate, 20, min(10, itemsHeight))
	l.FilterInput.Prompt = icon.Filter.Raw() + " "
	l.SetShowHelp(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowPagination(false)
	l.DisableQuitKeybindings()
	l.InfiniteScrolling = true

	l.Paginator.Type = paginator.Arabic

	l.SetStatusBarItemName(singular, plural)

	s := &Model{
		Model:    l,
		delegate: &delegate,
		KeyMap:   newKeyMap(&l.KeyMap),
	}
	s.updateKeybinds()
	return s
}
