package mangas

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/model/search"
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

// handleBrowsingCmd is the usual list behavior
func (s *state) handleBrowsingCmd(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering || s.search.State() == search.Searching {
			goto end
		}

		i, ok := s.list.SelectedItem().(*item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			if config.Download.Metadata.Search.Get() {
				return searchMetadataCmd(i)
			}
			return searchVolumesCmd(i)
		case key.Matches(msg, s.keyMap.search):
			s.list.ResetFilter()
			return s.search.Focus()
		}
	}
end:
	return s.list.Update(ctx, msg)
}
