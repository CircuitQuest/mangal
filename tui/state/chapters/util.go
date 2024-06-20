package chapters

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapsdownloaded"
	"github.com/luevano/mangal/tui/state/chapsdownloading"
	"github.com/zyedidia/generic/set"
)

func (s *state) readChapterCmd(ctx context.Context, path string, item *item, options libmangal.ReadOptions) tea.Cmd {
	chapter := item.chapter

	return tea.Sequence(
		actionRunningCmd("read"),
		base.Loading(fmt.Sprintf("Opening %q for reading", chapter)),
		func() tea.Msg {
			err := s.client.ReadChapter(ctx, path, chapter, options)
			if err != nil {
				return err
			}

			return nil
		},
		base.Loaded,
		actionRanCmd,
	)
}

func (s *state) downloadChapterCmd(ctx context.Context, item *item, options libmangal.DownloadOptions, readAfter bool) tea.Cmd {
	chapter := item.chapter

	if item.downloadedFormats.Has(options.Format) {
		return base.Notify(fmt.Sprintf("Chapter %q already downloaded in %s format", chapter, options.Format))
	}

	return tea.Sequence(
		actionRunningCmd("download"),
		base.Loading(fmt.Sprintf("Downloading %q", chapter)),
		func() tea.Msg {
			// TODO: make use of the returned data for data aggregation?
			downChap, err := s.client.DownloadChapter(ctx, chapter, options)
			if err != nil {
				return err
			}

			endCmd := base.Notify(fmt.Sprintf("Downloaded %q", chapter))
			if readAfter {
				endCmd = readChapterCmd(downChap.Path(), item, config.ReadOptions())
			}

			return tea.Sequence(updateItemCmd(item), endCmd)()
		},
		base.Loaded,
		actionRanCmd,
	)
}

// TODO: implement base.Loading/Loaded and actionRunningCmd/actionRanCmd
func (s *state) downloadChaptersCmd(items set.Set[*item], chapters []mangadata.Chapter, options libmangal.DownloadOptions) tea.Cmd {
	return func() tea.Msg {
		state := chapsdownloading.New(
			chapters,
			chapsdownloading.Options{
				DownloadChapter: func(ctx context.Context, chapter mangadata.Chapter) (*metadata.DownloadedChapter, error) {
					return s.client.DownloadChapter(ctx, chapter, options)
				},
				OnDownloadFinished: func(downChaps []*metadata.DownloadedChapter, succeed, failed []mangadata.Chapter) tea.Cmd {
					// TODO: better handle this, need to refactor chapsdownloading/chapsdownloaded
					// becasue the updateItemsCmd can't be send from here as it will be captured by the inner states
					defer s.updateItemsCmd(items)

					return func() tea.Msg {
						return chapsdownloaded.New(chapsdownloaded.Options{
							Succeed:          succeed,
							Failed:           failed,
							SucceedDownloads: downChaps,
							DownloadChapters: func(chapters []mangadata.Chapter) tea.Cmd {
								return s.downloadChaptersCmd(items, chapters, options)
							},
						})
					}
				},
			},
		)

		s.client.Logger().SetOnLog(func(format string, a ...any) {
			state.SetMessage(fmt.Sprintf(format, a...))
			log.Log(format, a...)
		})

		return state
	}
}
