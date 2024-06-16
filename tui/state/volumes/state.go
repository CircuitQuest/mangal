package volumes

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	list    *list.State
	volumes []mangadata.Volume
	manga   mangadata.Manga
	client  *libmangal.Client
	keyMap  keyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.list.FilterState() == _list.Unfiltered
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: s.manga.String()}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *State) Status() string {
	return s.list.Status()
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return s.list.Init(ctx)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		item, ok := s.list.SelectedItem().(Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Searching", fmt.Sprintf("Searching chapters for volume %s", item.volume))
				},
				func() tea.Msg {
					chapterList, err := s.client.VolumeChapters(ctx, item.volume)
					if err != nil {
						return err
					}

					return chapters.New(s.client, s.manga, item.volume, chapterList)
				},
			)
		}
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *State) View() string {
	return s.list.View()
}
