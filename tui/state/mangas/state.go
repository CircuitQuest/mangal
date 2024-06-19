package mangas

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
	"github.com/luevano/mangal/tui/state/volumes"
	list "github.com/luevano/mangal/tui/state/wrapper/list"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list   *list.State
	mangas []mangadata.Manga
	client *libmangal.Client
	query  string
	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.Backable()
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{Text: fmt.Sprintf("Search %q", s.query)}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *state) Status() string {
	return s.list.Status()
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
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
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
			return searchMetadataCmd(i)
		}
	case searchMetadataMsg:
		i := msg.item

		return tea.Sequence(
			base.Loading(fmt.Sprintf("Searching Anilist manga for %q", i.manga)),
			func() tea.Msg {
				// TODO: handle more cases for missing/partial metadata
				// Find anilist manga closest to the selected manga and assign it
				anilistManga, found, err := s.client.Anilist().SearchByManga(context.Background(), i.manga)
				if err != nil {
					return err
				}
				if !found {
					log.Log("Couldn't find Anilist for %q", i.manga)
				} else {
					i.manga.SetMetadata(anilistManga.Metadata())
					log.Log("Found and set Anilist for %q: %q (%d)", i.manga, anilistManga.String(), anilistManga.ID)
				}

				return searchVolumesMsg{i}
			},
			base.Loaded,
		)
	case searchVolumesMsg:
		i := msg.item

		return tea.Sequence(
			base.Loading(fmt.Sprintf("Searching volumes for %q", i.manga)),
			func() tea.Msg {
				volumeList, err := s.client.MangaVolumes(ctx, i.manga)
				if err != nil {
					return err
				}

				if len(volumeList) != 1 || !config.TUI.ExpandSingleVolume.Get() {
					return volumes.New(s.client, i.manga, volumeList)
				}

				// It's guaranteed to at least contain 1 volume
				volume := volumeList[0]
				chapterList, err := s.client.VolumeChapters(ctx, volume)
				if err != nil {
					return err
				}

				return chapters.New(s.client, i.manga, volume, chapterList)
			},
			base.Loaded,
		)
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
