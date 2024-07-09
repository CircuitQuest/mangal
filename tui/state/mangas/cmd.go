package mangas

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
)

func searchMangasCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return searchMangasMsg{
			query: query,
		}
	}
}

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
		if s.list.FilterState() == list.Filtering {
			goto end
		}

		i, ok := s.list.SelectedItem().(*item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			// TODO: only do a search metadata if the option is set?
			return searchMetadataCmd(i)
		case key.Matches(msg, s.keyMap.search):
			s.list.ResetFilter()
			s.searchState = searching
			s.searchInput.CursorEnd()
			s.searchInput.Focus()

			s.updateKeybindings()

			return textinput.Blink
		}
	}
end:
	return s.list.Update(ctx, msg)
}

// handleSearchingCmd controls the search bar behavior
func (s *state) handleSearchingCmd(_ context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.cancelSearch):
			s.searchState = searchCanceled
			s.searchInput.Blur()
			s.searchInput.Reset()
			s.updateKeybindings()

			// keep the last searched query
			s.searchInput.SetValue(s.query)

			return nil
		case key.Matches(msg, s.keyMap.confirmSearch):
			s.searchInput.Blur()
			s.searchState = searched
			s.updateKeybindings()

			s.query = s.searchInput.Value()
			return searchMangasCmd(s.query)
		}
	}
	input, inputUpdateCmd := s.searchInput.Update(msg)
	searchChanged := s.searchInput.Value() != input.Value()
	if searchChanged {
		s.keyMap.confirmSearch.SetEnabled(s.searchInput.Value() != "")
	}
	s.searchInput = input

	return inputUpdateCmd
}
