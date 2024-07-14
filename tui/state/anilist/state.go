package anilist

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/tui/model/search"
	metadataViewer "github.com/luevano/mangal/tui/state/metadata"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list    *list.Model
	search  *search.Model
	manga   mangadata.Manga
	anilist *lmanilist.Anilist

	searched bool

	keyMap *keyMap
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
		Text:       "Anilist Search",
		Background: color.Anilist,
		Foreground: color.Bright,
	}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	if !s.searched {
		return ""
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
		s.list.Init(),
	)
}

// Updateimplements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering || s.search.State() == search.Searching {
			goto end
		}

		// don't return on nil item, keybinds will be disabled for relevant actions
		i, _ := s.list.SelectedItem().(*item)
		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return s.setMetadataCmd(i.manga)
		case key.Matches(msg, s.keyMap.search):
			s.list.ResetFilter()
			return s.search.Focus()
		case key.Matches(msg, s.keyMap.metadata):
			return func() tea.Msg {
				return metadataViewer.New(&i.manga)
			}
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
	return s.list.Update(msg)
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

func (s *state) updateKeybinds() {
	enable := len(s.list.Items()) != 0
	s.keyMap.confirm.SetEnabled(enable)
	s.keyMap.metadata.SetEnabled(enable)
}
