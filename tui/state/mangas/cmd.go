package mangas

import tea "github.com/charmbracelet/bubbletea"

func searchMetadataCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return searchMetadataMsg{
			item: item,
		}
	}
}

func searchVolumesCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return searchVolumesMsg{
			item: item,
		}
	}
}
