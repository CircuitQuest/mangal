package chapters

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapsdownloaded"
	"github.com/luevano/mangal/tui/state/chapsdownloading"
	"github.com/luevano/mangal/tui/state/loading"
)

func (s *State) downloadChaptersCmd(chapters []libmangal.Chapter, options libmangal.DownloadOptions) tea.Cmd {
	return func() tea.Msg {
		state := chapsdownloading.New(
			chapters,
			chapsdownloading.Options{
				DownloadChapter: func(ctx context.Context, chapter libmangal.Chapter) (libmangal.DownloadedChapter, error) {
					return s.client.DownloadChapter(ctx, chapter, options)
				},
				OnDownloadFinished: func(downChaps []libmangal.DownloadedChapter, succeed, failed []libmangal.Chapter) tea.Cmd {
					return func() tea.Msg {
						return chapsdownloaded.New(chapsdownloaded.Options{
							Succeed:          succeed,
							Failed:           failed,
							SucceedDownloads: downChaps,
							DownloadChapters: func(chapters []libmangal.Chapter) tea.Cmd {
								return s.downloadChaptersCmd(chapters, options)
							},
						})
					}
				},
			},
		)

		s.client.Logger().SetOnLog(func(msg string) {
			state.SetMessage(msg)
			log.Log(msg)
		})

		return state
	}
}

func (s *State) downloadChapterCmd(ctx context.Context, chapter libmangal.Chapter, options libmangal.DownloadOptions) tea.Cmd {
	volume := chapter.Volume()
	manga := volume.Manga()

	loadingState := loading.New("Downloading", fmt.Sprintf("%s / Vol. %s / %s", manga, volume, chapter))
	return tea.Sequence(
		func() tea.Msg {
			return loadingState
		},
		func() tea.Msg {
			s.client.Logger().SetOnLog(func(msg string) {
				loadingState.SetMessage(msg)
				log.Log(msg)
			})

			// TODO: make use of the returned data for data aggregation?
			_, err := s.client.DownloadChapter(ctx, chapter, options)
			if err != nil {
				return err
			}

			return base.MsgBack{}
		},
	)
}
