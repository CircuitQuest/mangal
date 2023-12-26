package mangas

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/volumes"
	"github.com/luevano/libmangal"
)

var _ base.State = (*State)(nil)

type State struct {
	query  string
	client *libmangal.Client
	mangas []libmangal.Manga
	list   *listwrapper.State
	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: fmt.Sprintf("Search %q", s.query)}
}

func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

func (s *State) Status() string {
	return s.list.Status()
}

func (s *State) Backable() bool {
	return s.list.Backable()
}

func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			goto end
		}

		item, ok := s.list.SelectedItem().(Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.Confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Loading...", fmt.Sprintf("Getting volumes for %q", item.Manga))
				},
				func() tea.Msg {
					v, err := s.client.MangaVolumes(model.Context(), item.Manga)
					if err != nil {
						return err
					}

					if len(v) != 1 || !config.Config.TUI.ExpandSingleVolume.Get() {
						return volumes.New(s.client, item.Manga, v)
					}

					volume := v[0]
					c, err := s.client.VolumeChapters(model.Context(), volume)
					if err != nil {
						return err
					}

					return chapters.New(s.client, volume, c)
				},
			)
		}
	}
end:
	return s.list.Update(model, msg)
}

func (s *State) View(model base.Model) string {
	return s.list.View(model)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return s.list.Init(model)
}
