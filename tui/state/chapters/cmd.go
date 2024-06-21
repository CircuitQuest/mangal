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

func (s *state) blockedActionByCmd(wanted string) tea.Cmd {
	return base.Notify(fmt.Sprintf("Can't perform %q right now, %q is running", wanted, s.actionRunning))
}

func (s *state) readChapterCmd(ctx context.Context, path string, item *item, options libmangal.ReadOptions) tea.Cmd {
	chapter := item.chapter

	return tea.Sequence(
		base.Loading(fmt.Sprintf("Opening %q for reading", chapter)),
		func() tea.Msg {
			s.actionRunningNow("read")
			defer s.actionRunningNow("")

			err := s.client.ReadChapter(ctx, path, chapter, options)
			if err != nil {
				return err
			}

			return nil
		},
		base.Loaded,
	)
}

func (s *state) downloadChapterCmd(ctx context.Context, item *item, options libmangal.DownloadOptions, readAfter bool) tea.Cmd {
	chapter := item.chapter

	if item.downloadedFormats.Has(options.Format) {
		return base.Notify(fmt.Sprintf("Chapter %q already downloaded in %s format", chapter, options.Format))
	}

	return tea.Sequence(
		base.Loading(fmt.Sprintf("Downloading %q", chapter)),
		func() tea.Msg {
			s.actionRunningNow("download")
			defer s.actionRunningNow("")

			// TODO: make use of the returned data for data aggregation?
			downChap, err := s.client.DownloadChapter(ctx, chapter, options)
			if err != nil {
				return err
			}
			s.updateItem(item)

			if readAfter {
				return readChapterMsg{downChap.Path(), item, config.ReadOptions()}
			}
			return base.Notify(fmt.Sprintf("Downloaded %q", chapter))()
		},
		base.Loaded,
	)
}

// TODO: temporary placeholder that just wraps the msg, the msg will trigger the actual cmd (below)
func downloadChaptersCmd(items set.Set[*item], options libmangal.DownloadOptions) tea.Cmd {
	return func() tea.Msg {
		return downloadChaptersMsg{
			items:   items,
			options: options,
		}
	}
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
					s.updateItems(items)

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
