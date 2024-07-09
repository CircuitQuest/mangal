package anilistmangas

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func searchMangasCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return searchMangasMsg{
			query: query,
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

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			i, ok := s.list.SelectedItem().(*item)
			if !ok {
				return nil
			}

			return s.onResponse(i.manga)
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
