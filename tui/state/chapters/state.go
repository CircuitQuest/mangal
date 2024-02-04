package chapters

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/anilistmangas"
	"github.com/luevano/mangal/tui/state/confirm"
	"github.com/luevano/mangal/tui/state/formats"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/state/loading"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/skratchdot/open-golang/open"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*State)(nil)

type State struct {
	client            *libmangal.Client
	volume            libmangal.Volume
	selected          set.Set[*Item]
	list              *listwrapper.State
	keyMap            KeyMap
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool
}

func (s *State) Intermediate() bool {
	return false
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	volume := s.volume
	manga := volume.Manga()

	return base.Title{Text: fmt.Sprintf("%s / Vol. %d", manga, volume.Info().Number)}
}

func (s *State) Subtitle() string {
	var subtitle strings.Builder

	subtitle.WriteString(s.list.Subtitle())

	if s.selected.Size() > 0 {
		subtitle.WriteString(" ")
		subtitle.WriteString(fmt.Sprint(s.selected.Size()))
		subtitle.WriteString(" selected")
	}

	subtitle.WriteString(". Download ")
	subtitle.WriteString(config.Config.Download.Format.Get().String())
	subtitle.WriteString(" & Read ")
	subtitle.WriteString(config.Config.Read.Format.Get().String())

	return subtitle.String()
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

		item, ok := s.list.SelectedItem().(*Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.Toggle):
			item.Toggle()

			return nil
		case key.Matches(msg, s.keyMap.UnselectAll):
			for _, item := range s.selected.Keys() {
				item.Toggle()
			}

			return nil
		case key.Matches(msg, s.keyMap.SelectAll):
			for _, listItem := range s.list.Items() {
				item, ok := listItem.(*Item)
				if !ok {
					continue
				}

				if !item.IsSelected() {
					item.Toggle()
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.ChangeFormat):
			return func() tea.Msg {
				return formats.New()
			}
		case key.Matches(msg, s.keyMap.OpenURL):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Opening", item.chapter.String())
				},
				func() tea.Msg {
					err := open.Run(item.chapter.Info().URL)
					if err != nil {
						return err
					}

					return base.MsgBack{}
				},
			)
		case key.Matches(msg, s.keyMap.Download) || (s.selected.Size() > 0 && key.Matches(msg, s.keyMap.Confirm)):
			var chapters []libmangal.Chapter

			if s.selected.Size() == 0 {
				chapters = append(chapters, item.chapter)
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

						return s.downloadChaptersCmd(chapters, config.Config.DownloadOptions())
					},
				)
			}
		case key.Matches(msg, s.keyMap.Read) || (s.selected.Size() == 0 && key.Matches(msg, s.keyMap.Confirm)):
			// If download on read is wanted, then use the normal download path
			var directory string
			if config.Config.Read.DownloadOnRead.Get() {
				directory = path.DownloadsDir()
			} else {
				directory = path.TempDir()
			}

			// Modify a bit the configured download options for this
			downloadOptions := config.Config.DownloadOptions()
			downloadOptions.Format = config.Config.Read.Format.Get()
			downloadOptions.Directory = directory
			downloadOptions.SkipIfExists = true
			downloadOptions.ReadAfter = true
			downloadOptions.CreateProviderDir = true
			downloadOptions.CreateMangaDir = true
			downloadOptions.CreateVolumeDir = true

			if item.DownloadedFormats().Has(downloadOptions.Format) {
				return tea.Sequence(
					func() tea.Msg {
						return loading.New("Opening for reading", item.chapter.String())
					},
					func() tea.Msg {
						err := s.client.ReadChapter(
							model.Context(),
							item.Path(downloadOptions.Format),
							item.chapter,
							downloadOptions.ReadOptions,
						)
						if err != nil {
							return err
						}

						return base.MsgBack{}
					},
				)
			}

			return s.downloadChapterCmd(model.Context(), item.chapter, downloadOptions)
		// TODO: this should be set some levels before
		case key.Matches(msg, s.keyMap.Anilist):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Searching", "Getting Anilist Mangas")
				},
				func() tea.Msg {
					var mangas []libmangal.AnilistManga

					// TODO: revert to just Title instead of AnilistSearch?
					var mangaTitle string
					mangaInfo := item.chapter.Volume().Manga().Info()
					if mangaInfo.AnilistSearch != "" {
						mangaTitle = mangaInfo.AnilistSearch
					} else {
						mangaTitle = mangaInfo.Title
					}

					closest, ok, err := s.client.Anilist().FindClosestManga(model.Context(), mangaTitle)
					if err != nil {
						return err
					}

					if ok {
						mangas = append(mangas, closest)
					}

					mangaSearchResults, err := s.client.Anilist().SearchMangas(model.Context(), mangaTitle)
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
						func(response *libmangal.AnilistManga) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									err := s.client.Anilist().BindTitleWithID(mangaTitle, response.ID)
									if err != nil {
										return err
									}
									// TODO: need to set the AnilistManga to the root Manga...
									// s.volume.Manga().SetAnilistManga(*response)

									return base.MsgBack{}
								},
								s.list.Notify("Binded to "+response.Title.English, time.Second*3),
							)
						},
					)
				},
			)
		case key.Matches(msg, s.keyMap.ToggleChapterNumber):
			*s.showChapterNumber = !(*s.showChapterNumber)

			return s.list.Update(model, msg)
		case key.Matches(msg, s.keyMap.ToggleGroup):
			*s.showGroup = !(*s.showGroup)

			return s.list.Update(model, msg)
		case key.Matches(msg, s.keyMap.ToggleDate):
			*s.showDate = !(*s.showDate)

			return s.list.Update(model, msg)
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
