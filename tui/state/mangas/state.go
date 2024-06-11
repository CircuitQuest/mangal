package mangas

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/volumes"
)

var _ base.State = (*State)(nil)

type State struct {
	query  string
	client *libmangal.Client
	mangas []*mangadata.Manga
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
					return loading.New("Searching", fmt.Sprintf("Finding Anilist for %q", *item.manga))
				},
				func() tea.Msg {
					// TODO: handle more cases for missing/partial metadata
					// Find anilist manga closest to the selected manga and assign it
					anilistManga, found, err := s.client.Anilist().SearchByManga(context.Background(), *item.manga)
					if err != nil {
						return err
					}
					if !found {
						log.Log("Couldn't find Anilist for %q", *item.manga)
					} else {
						(*item.manga).SetMetadata(anilistManga.Metadata())
						log.Log("Found and set Anilist for %q: %q (%d)", *item.manga, anilistManga.String(), anilistManga.ID)
					}

					return nil
				},
				func() tea.Msg {
					return loading.New("Searching", fmt.Sprintf("Getting volumes for %q", *item.manga))
				},
				func() tea.Msg {
					vL, err := s.client.MangaVolumes(model.Context(), *item.manga)
					if err != nil {
						return err
					}

					var volumeList []*mangadata.Volume
					for _, v := range vL {
						volumeList = append(volumeList, &v)
					}

					if len(vL) != 1 || !config.TUI.ExpandSingleVolume.Get() {
						return volumes.New(s.client, item.manga, volumeList)
					}

					volume := volumeList[0]
					cL, err := s.client.VolumeChapters(model.Context(), *volume)
					if err != nil {
						return err
					}

					var chapterList []*mangadata.Chapter
					for _, c := range cL {
						chapterList = append(chapterList, &c)
					}

					return chapters.New(s.client, item.manga, volume, chapterList)
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
