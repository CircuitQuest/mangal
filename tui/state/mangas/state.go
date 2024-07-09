package mangas

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
	"github.com/luevano/mangal/tui/state/volumes"
	list "github.com/luevano/mangal/tui/state/wrapper/list"
)

type searchState int

const (
	unsearched searchState = iota
	searching
	searched
	searchCanceled
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list   *list.State
	client *libmangal.Client

	query       string
	searchInput textinput.Model
	searchState searchState

	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.Backable() && s.searchState != searching
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
	if s.searchState == unsearched {
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
	s.searchInput.Width = size.Width
	sSize := base.Size{}
	sSize.Width, sSize.Height = lipgloss.Size(s.searchView())

	final := size
	final.Height -= sSize.Height

	return s.list.Resize(final)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	s.searchInput.CursorEnd()

	return tea.Sequence(
		s.searchInput.Focus(),
		textinput.Blink,
		s.list.Init(ctx),
	)
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case searchMangasMsg:
		query := msg.query

		return tea.Sequence(
			base.Loading(fmt.Sprintf("Searching for %q", query)),
			func() tea.Msg {
				mangas, err := s.client.SearchMangas(ctx, query)
				if err != nil {
					return nil
				}
				items := make([]_list.Item, len(mangas))

				for i, m := range mangas {
					items[i] = &item{manga: m}
				}

				s.list.SetItems(items)
				return nil
			},
			base.Loaded,
		)
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
				vols := len(volumeList)

				if config.TUI.ExpandAllVolumes.Get() {
					return searchAllChaptersMsg{i.manga, volumeList}
				}

				if vols == 1 && config.TUI.ExpandSingleVolume.Get() {
					return searchChaptersMsg{i.manga, volumeList[0]}
				}

				return volumes.New(s.client, i.manga, volumeList)
			},
			base.Loaded,
		)
	case searchAllChaptersMsg:
		manga := msg.manga
		volumes := msg.volumes

		// TODO: make different loading messages for each volume?
		return tea.Sequence(
			base.NotifyWithDuration(fmt.Sprintf("Skipped selecting volumes (cfg: %s)", config.TUI.ExpandAllVolumes.Key), 3*time.Second),
			base.Loading("Searching chapters for all volumes"),
			func() tea.Msg {
				var chapterList []mangadata.Chapter
				for _, v := range volumes {
					c, err := s.client.VolumeChapters(ctx, v)
					if err != nil {
						return err
					}
					chapterList = append(chapterList, c...)
				}

				return chapters.New(s.client, manga, nil, chapterList)
			},
			base.Loaded,
		)
	case searchChaptersMsg:
		manga := msg.manga
		volume := msg.volume

		return tea.Sequence(
			base.NotifyWithDuration(fmt.Sprintf("Skipped single volume (cfg: %s)", config.TUI.ExpandSingleVolume.Key), 3*time.Second),
			base.Loading("Searching chapters"),
			func() tea.Msg {
				chapterList, err := s.client.VolumeChapters(ctx, volume)
				if err != nil {
					return err
				}

				return chapters.New(s.client, manga, nil, chapterList)
			},
			base.Loaded,
		)
	}

	if s.searchState == searching || s.searchState == unsearched {
		return s.handleSearchingCmd(ctx, msg)
	}
	return s.handleBrowsingCmd(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	if s.searchState == unsearched {
		return s.searchView()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.searchView(),
		s.list.View(),
	)
}
