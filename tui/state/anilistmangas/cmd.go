package anilistmangas

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/model/search"
)

// handleBrowsingCmd is the usual list behavior
func (s *state) handleBrowsingCmd(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering || s.search.State() == search.Searching {
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
			return s.search.Focus()
		}
	}
end:
	return s.list.Update(ctx, msg)
}
