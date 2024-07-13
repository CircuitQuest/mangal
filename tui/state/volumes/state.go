package volumes

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/metadata"
	"github.com/luevano/mangal/tui/state/anilist"
	metadataViewer "github.com/luevano/mangal/tui/state/metadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list *list.State
	meta *metadata.Model

	volumes []mangadata.Volume
	manga   mangadata.Manga
	client  *libmangal.Client

	keyMap *keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.FilterState() == _list.Unfiltered
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{Text: s.manga.String()}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *state) Status() string {
	return s.meta.View() + " " + s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	return s.list.Init(ctx)
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		i, ok := s.list.SelectedItem().(*item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return s.searchVolumeChapters(ctx, i)
		case key.Matches(msg, s.keyMap.anilist):
			return func() tea.Msg {
				return anilist.New(s.client.Anilist(), s.manga)
			}
		case key.Matches(msg, s.keyMap.metadata):
			return func() tea.Msg {
				return metadataViewer.New(s.manga.Metadata())
			}
		case key.Matches(msg, s.keyMap.info):
			s.meta.ShowFull = !s.meta.ShowFull
		}
	case base.RestoredMsg:
		// in case the metadata was updated in the anilist state
		s.meta.SetMetadata(s.manga.Metadata())
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
