package chapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/anilistmangas"
	"github.com/luevano/mangal/tui/state/confirm"
	"github.com/luevano/mangal/tui/state/formats"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list *list.State

	chapters []mangadata.Chapter
	volume   mangadata.Volume // can be nil
	manga    mangadata.Manga
	client   *libmangal.Client

	selected set.Set[*item]

	renderedSep             string
	renderedSubtitleFormats string

	// to avoid extra read/download cmds from
	// firing up when an action is already happening,
	// only blocks keymaps handled by Update for read/download
	actionRunning string

	showVolumeNumber  *bool
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool

	styles styles
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
	return base.Title{Text: s.manga.String()}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	var subtitle strings.Builder
	subtitle.Grow(100)

	subtitle.WriteString(s.list.Subtitle())
	if s.selected.Size() > 0 {
		selected := s.renderedSep +
			s.styles.subtitle.Render(fmt.Sprintf("%d selected", s.selected.Size()))
		subtitle.WriteString(selected)
	}
	subtitle.WriteString(s.renderedSubtitleFormats)

	return subtitle.String()
}

// Status implements base.State.
func (s *state) Status() string {
	if s.volume != nil {
		return fmt.Sprintf("Vol. %s%s%s", s.volume, s.renderedSep, s.list.Status())
	}
	return s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	s.updateRenderedSubtitleFormats()

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
		case key.Matches(msg, s.keyMap.toggle):
			i.toggle()
			if i.selected {
				s.selected.Put(i)
			} else {
				s.selected.Remove(i)
			}

			return nil
		case key.Matches(msg, s.keyMap.unselectAll):
			for _, item := range s.selected.Keys() {
				item.toggle()
				s.selected.Remove(item)
			}

			return nil
		case key.Matches(msg, s.keyMap.selectAll):
			for _, listItem := range s.list.Items() {
				it, ok := listItem.(*item)
				if !ok {
					continue
				}

				if !it.selected {
					it.toggle()
					s.selected.Put(it)
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.changeFormat):
			return func() tea.Msg {
				// all this does is actually change the config for the formats
				return formats.New()
			}
		case key.Matches(msg, s.keyMap.openURL):
			return s.openURLCmd(i)
		case key.Matches(msg, s.keyMap.download):
			if s.actionRunning != "" {
				return s.blockedActionByCmd("download")
			}

			// when no toggled chapters then just download the one hovered
			if s.selected.Size() == 0 {
				// TODO: add confirmation?
				return s.downloadChapterCmd(ctx, i, config.DownloadOptions(), false)
			}

			// TODO: refactor confirmation state?
			return func() tea.Msg {
				return confirm.New(
					fmt.Sprint("Download ", stringutil.Quantify(s.selected.Size(), "chapter", "chapters")),
					func(response bool) tea.Cmd {
						if !response {
							return base.Back
						}

						return s.downloadChaptersCmd(s.selected, config.DownloadOptions())
					},
				)
			}
		case key.Matches(msg, s.keyMap.read):
			if s.actionRunning != "" {
				return s.blockedActionByCmd("read")
			}

			// when no toggled chapters then just download the one selected
			if s.selected.Size() > 1 {
				return base.Notify("Can't open for reading more than 1 chapter")
			}

			// use the toggled item, else the hovered one
			it := i
			if s.selected.Size() == 1 {
				it = s.selected.Keys()[0]
			}

			if it.readAvailablePath != "" {
				log.Log("Read format already downloaded")
				return s.readChapterCmd(ctx, it.readAvailablePath, it, config.ReadOptions())
			}

			downloadOptions := config.DownloadOptions()
			// TODO: add warning when read format != download format?
			downloadOptions.Format = config.Read.Format.Get()
			// If shouldn't download on read, save to tmp dir with all dirs created
			if !config.Read.DownloadOnRead.Get() {
				downloadOptions.Directory = path.TempDir()
				downloadOptions.CreateProviderDir = true
				downloadOptions.CreateMangaDir = true
				downloadOptions.CreateVolumeDir = true
			}

			// TODO: add confirmation?
			log.Log("Read format not yet downloaded, downloading")
			return s.downloadChapterCmd(ctx, it, downloadOptions, true)
		// TODO: refactor/fix this so that the metadata is propagated (probably needs a change on libmangal itself)
		case key.Matches(msg, s.keyMap.anilist):
			return tea.Sequence(
				base.Loading(fmt.Sprintf("Searching Anilist mangas for %q", s.manga)),
				func() tea.Msg {
					var mangas []lmanilist.Manga

					// TODO: solidify the metadata gathering, missing/partial
					// TODO: revert to just Title instead of AnilistSearch?
					var mangaTitle string
					mangaInfo := i.chapter.Volume().Manga().Info()
					if mangaInfo.AnilistSearch != "" {
						mangaTitle = mangaInfo.AnilistSearch
					} else {
						mangaTitle = mangaInfo.Title
					}

					closest, ok, err := s.client.Anilist().FindClosestManga(ctx, mangaTitle)
					if err != nil {
						return err
					}

					if ok {
						mangas = append(mangas, closest)
					}

					mangaSearchResults, err := s.client.Anilist().SearchMangas(ctx, mangaTitle)
					if err != nil {
						return nil
					}

					for _, manga := range mangaSearchResults {
						if manga.ID == closest.ID {
							continue
						}

						mangas = append(mangas, manga)
					}

					return anilistmangas.New(
						s.client.Anilist(),
						mangas,
						func(response *lmanilist.Manga) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									log.Log("Setting Anilist manga %q (%d)", response.String(), response.ID)
									s.manga.SetMetadata(response.Metadata())

									return nil
								},
								base.NotifyWithDuration(fmt.Sprintf("Set Anilist %s (%d)", response.String(), response.ID), 3*time.Second),
							)
						},
					)
				},
				base.Loaded,
			)
		case key.Matches(msg, s.keyMap.toggleVolumeNumber):
			*s.showVolumeNumber = !(*s.showVolumeNumber)
		case key.Matches(msg, s.keyMap.toggleChapterNumber):
			*s.showChapterNumber = !(*s.showChapterNumber)
		case key.Matches(msg, s.keyMap.toggleGroup):
			*s.showGroup = !(*s.showGroup)
			s.updateListDelegate()
		case key.Matches(msg, s.keyMap.toggleDate):
			*s.showDate = !(*s.showDate)
			s.updateListDelegate()
		}
	case readChapterMsg:
		return s.readChapterCmd(ctx, msg.path, msg.item, msg.options)
	case downloadChapterMsg:
		return s.downloadChapterCmd(ctx, msg.item, msg.options, msg.readAfter)
	case downloadChaptersMsg:
		return s.downloadChaptersCmd(msg.items, msg.options)
	case base.RestoredMsg:
		// update the items sent for download when coming back
		s.updateItems(s.selected)
		// can't distinguish from which state this is being restored
		s.updateRenderedSubtitleFormats()
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
