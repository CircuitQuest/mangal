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
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/anilist"
	metadataViewer "github.com/luevano/mangal/tui/state/metadata"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list   *list.Model
	search *search.Model

	client *libmangal.Client

	searched bool

	extraInfo     *bool
	fullExtraInfo *bool

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
		s.list.Init(),
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
		case key.Matches(msg, s.keyMap.anilist):
			return func() tea.Msg {
				return anilist.New(s.client.Anilist(), i.manga)
			}
		case key.Matches(msg, s.keyMap.metadata):
			return func() tea.Msg {
				return metadataViewer.New(i.manga.Metadata())
			}
		// TODO: only toggle for hovered/selected item? (both info and full info)
		case key.Matches(msg, s.keyMap.info):
			*s.extraInfo = !(*s.extraInfo)
			s.updateKeybinds() // to enable/disable toggleFullMetadata kb
		case key.Matches(msg, s.keyMap.toggleFullMeta):
			*s.fullExtraInfo = !(*s.fullExtraInfo)
		}
	case search.SearchMsg:
		return s.searchMangasCmd(ctx, string(msg))
	case search.SearchCancelMsg:
		if s.search.Query() == "" {
			return base.Back
		}
	case base.RestoredMsg:
		// in case that the metadata was updated, update all items
		s.updateAllItems()
	}
end:
	if s.search.State() == search.Searching {
		return s.search.Update(msg)
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
