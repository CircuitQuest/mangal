package chapters

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/zyedidia/generic/set"
)

func actionRunningCmd(action string) tea.Cmd {
	return func() tea.Msg {
		return actionRunningMsg{
			action: action,
		}
	}
}

func actionRanCmd() tea.Msg {
	return actionRunningMsg{
		action: "",
	}
}

func blockedActionCmd(wanted string) tea.Cmd {
	return func() tea.Msg {
		return blockedActionMsg{
			wanted: wanted,
		}
	}
}

func updateItemCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		return updateItemMsg{
			item: item,
		}
	}
}

func updateItemsCmd(items set.Set[*item]) tea.Cmd {
	return func() tea.Msg {
		return updateItemsMsg{
			items: items,
		}
	}
}

func readChapterCmd(path string, item *item, options libmangal.ReadOptions) tea.Cmd {
	return func() tea.Msg {
		return readChapterMsg{
			path:    path,
			item:    item,
			options: options,
		}
	}
}

func downloadChapterCmd(item *item, options libmangal.DownloadOptions, readAfter bool) tea.Cmd {
	return func() tea.Msg {
		return downloadChapterMsg{
			item:      item,
			options:   options,
			readAfter: readAfter,
		}
	}
}

func downloadChaptersCmd(items set.Set[*item], options libmangal.DownloadOptions) tea.Cmd {
	return func() tea.Msg {
		return downloadChaptersMsg{
			items:   items,
			options: options,
		}
	}
}
