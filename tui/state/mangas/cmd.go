package mangas

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
)

func searchMetadataCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		return searchMetadataMsg{
			item: item,
		}
	}
}

func searchVolumesCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		return searchVolumesMsg{
			item: item,
		}
	}
}

func searchChaptersCmd(manga mangadata.Manga, volume mangadata.Volume) tea.Cmd {
	return func() tea.Msg {
		return searchChaptersMsg{
			manga:  manga,
			volume: volume,
		}
	}
}

func searchAllChaptersCmd(manga mangadata.Manga, volumes []mangadata.Volume) tea.Cmd {
	return func() tea.Msg {
		return searchAllChaptersMsg{
			manga:   manga,
			volumes: volumes,
		}
	}
}
