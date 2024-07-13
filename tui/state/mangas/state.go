package mangas

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list   *list.State
	search *search.Model
	client *libmangal.Client

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
	return base.Title{Text: s.client.Info().Name}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	if !s.searched {
		return "Search on " + s.client.Info().Name
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
		s.search.Focus(),
		s.list.Init(ctx),
	)
}

// Update implements base.State.
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
			if config.Download.Metadata.Search.Get() {
				return s.searchMetadataCmd(ctx, i)
			}
			return s.searchVolumesCmd(ctx, i)
		case key.Matches(msg, s.keyMap.search):
			s.list.ResetFilter()
			return s.search.Focus()
		}
	case search.SearchMsg:
		return s.searchMangasCmd(ctx, string(msg))
	case search.SearchCancelMsg:
		if s.search.Query() == "" {
			return base.Back
		}
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

func (s *state) updateKeybinds() {
	s.keyMap.confirm.SetEnabled(len(s.list.Items()) != 0)
}
