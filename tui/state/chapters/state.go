package chapters

import (
	"context"
	"fmt"
	"sort"
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
	"github.com/skratchdot/open-golang/open"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list              *list.State
	chapters          []mangadata.Chapter
	volume            mangadata.Volume // can be nil
	manga             mangadata.Manga
	client            *libmangal.Client
	selected          *set.Set[*item]
	keyMap            keyMap
	showVolumeNumber  *bool
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool
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

	subtitle.WriteString(s.list.Subtitle())

	if s.selected.Size() > 0 {
		subtitle.WriteString(" ")
		subtitle.WriteString(fmt.Sprint(s.selected.Size()))
		subtitle.WriteString(" selected")
	}

	subtitle.WriteString(". Download ")
	subtitle.WriteString(config.Download.Format.Get().String())
	subtitle.WriteString(" & Read ")
	subtitle.WriteString(config.Read.Format.Get().String())

	return subtitle.String()
}

// Status implements base.State.
func (s *state) Status() string {
	if s.volume != nil {
		return fmt.Sprintf("Vol. %s%s%s", s.volume, base.StatusSeparator, s.list.Status())
	}
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
		case key.Matches(msg, s.keyMap.toggle):
			i.toggle()

			return nil
		case key.Matches(msg, s.keyMap.unselectAll):
			for _, item := range s.selected.Keys() {
				item.toggle()
			}

			return nil
		case key.Matches(msg, s.keyMap.selectAll):
			for _, listItem := range s.list.Items() {
				it, ok := listItem.(*item)
				if !ok {
					continue
				}

				if !it.isSelected() {
					it.toggle()
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.changeFormat):
			return func() tea.Msg {
				return formats.New()
			}
		case key.Matches(msg, s.keyMap.openURL):
			return tea.Sequence(
				base.Loading(fmt.Sprintf("Opening URL %s for chapter %q", i.chapter.Info().URL, i.chapter)),
				func() tea.Msg {
					err := open.Run(i.chapter.Info().URL)
					if err != nil {
						return err
					}

					return nil
				},
				base.Loaded,
			)
		case key.Matches(msg, s.keyMap.download):
			var chapters []mangadata.Chapter

			if s.selected.Size() == 0 {
				chapters = append(chapters, i.chapter)
			} else {
				for _, item := range s.selected.Keys() {
					chapters = append(chapters, item.chapter)
				}
			}

			sort.SliceStable(chapters, func(i, j int) bool {
				return chapters[i].Info().Number < chapters[j].Info().Number
			})

			return func() tea.Msg {
				return confirm.New(
					fmt.Sprint("Download ", stringutil.Quantify(len(chapters), "chapter", "chapters")),
					func(response bool) tea.Cmd {
						if !response {
							return base.Back
						}

						return s.downloadChaptersCmd(chapters, config.DownloadOptions())
					},
				)
			}
		case key.Matches(msg, s.keyMap.read):
			// If download on read is wanted, then use the normal download path
			var directory string
			if config.Read.DownloadOnRead.Get() {
				directory = config.Download.Path.Get()
			} else {
				directory = path.TempDir()
			}

			// Modify a bit the configured download options for this
			downloadOptions := config.DownloadOptions()
			downloadOptions.Format = config.Read.Format.Get()
			downloadOptions.Directory = directory
			downloadOptions.SkipIfExists = true
			downloadOptions.CreateProviderDir = true
			downloadOptions.CreateMangaDir = true
			downloadOptions.CreateVolumeDir = true

			if i.downloadedFormats().Has(downloadOptions.Format) {
				return tea.Sequence(
					base.Loading(fmt.Sprintf("Opening %q for reading", i.chapter)),
					func() tea.Msg {
						err := s.client.ReadChapter(
							ctx,
							i.path(downloadOptions.Format),
							i.chapter,
							config.ReadOptions(),
						)
						if err != nil {
							return err
						}

						return nil
					},
					base.Loaded,
				)
			}

			// TODO: now that libmangal doesn't have an option to "read after downloading",
			// handle that case by running read chapter again (like above) after the download
			return s.downloadChapterCmd(ctx, i.chapter, downloadOptions)
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
		case key.Matches(msg, s.keyMap.toggleDate):
			*s.showDate = !(*s.showDate)
		}
	}

end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
