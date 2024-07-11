package anilistmangas

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

type onResponseFunc func(manga lmanilist.Manga) tea.Cmd

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list    *list.State
	search  *search.Model
	anilist *lmanilist.Anilist

	searched bool

	onResponse onResponseFunc

	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.Backable() && s.search.State() != search.Searching
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{
		Text:       "Anilist Mangas",
		Background: color.Anilist.Background,
		Foreground: color.Anilist.Foreground,
	}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	if !s.searched {
		return "Search on Anilist"
	}
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *state) Status() string {
	return s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	s.search.Resize(size)
	_, searchHeight := lipgloss.Size(s.search.View())
	size.Height -= searchHeight

	return s.list.Resize(size)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	return tea.Sequence(
		s.searchCmd(ctx, s.search.Query()),
		s.list.Init(ctx),
	)
}

// Updateimplements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering || s.search.State() == search.Searching {
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
	case search.SearchMsg:
		return s.searchCmd(ctx, string(msg))
	}
end:
	if s.search.State() == search.Searching {
		input, updateCmd := s.search.Update(msg)
		s.search = input.(*search.Model)
		return updateCmd
	}
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	if !s.searched {
		return s.search.View()
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.search.View(),
		s.list.View(),
	)
}
