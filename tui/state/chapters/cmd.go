package chapters

import (
	"context"
	"fmt"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata/anilist"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/anilistmangas"
	"github.com/luevano/mangal/tui/state/confirm"
	"github.com/luevano/mangal/tui/state/download"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/skratchdot/open-golang/open"
	"github.com/zyedidia/generic/set"
)

func (s *state) updateMetadataCmd(anilistManga anilist.Manga) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			log.Log("Setting Anilist %q (%d)", anilistManga.String(), anilistManga.ID)
			s.manga.SetMetadata(anilistManga.Metadata())

			return nil
		},
		base.NotifyWithDuration(fmt.Sprintf("Set Anilist %s (%d)", anilistManga.String(), anilistManga.ID), 3*time.Second),
	)
}

func (s *state) blockedActionByCmd(wanted string) tea.Cmd {
	return base.Notify(fmt.Sprintf("Can't perform %q right now, %q is running", wanted, s.actionRunning))
}

func (s *state) openURLCmd(chapter mangadata.Chapter) tea.Cmd {
	return tea.Sequence(
		base.Loading(fmt.Sprintf("Opening URL %s for chapter %q", chapter.Info().URL, chapter)),
		func() tea.Msg {
			err := open.Run(chapter.Info().URL)
			if err != nil {
				return err
			}

			return nil
		},
		base.Loaded,
	)
}

func (s *state) downloadCmd(ctx context.Context, item *item) tea.Cmd {
	if s.actionRunning != "" {
		return s.blockedActionByCmd("download")
	}

	// when no toggled chapters then just download the one hovered
	if s.selected.Size() == 0 {
		// TODO: add confirmation?
		return s.downloadChapterCmd(ctx, item, config.DownloadOptions(), false)
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
				return s.readChapterCmd(ctx, downChap.Path(), item, config.ReadOptions())()
			}
			return base.Notify(fmt.Sprintf("Downloaded %q", chapter))()
		},
		base.Loaded,
	)
}

// TODO: implement base.Loading/Loaded and actionRunningCmd/actionRanCmd
func (s *state) downloadChaptersCmd(items set.Set[*item], options libmangal.DownloadOptions) tea.Cmd {
	return func() tea.Msg {
		var chapters []mangadata.Chapter
		for _, item := range items.Keys() {
			chapters = append(chapters, item.chapter)
		}
		sort.SliceStable(chapters, func(i, j int) bool {
			return chapters[i].Info().Number < chapters[j].Info().Number
		})

		return download.New(
			s.client,
			chapters,
			options,
		)
	}
}

func (s *state) readCmd(ctx context.Context, item *item) tea.Cmd {
	if s.actionRunning != "" {
		return s.blockedActionByCmd("read")
	}

	// when no toggled chapters then just download the one selected
	if s.selected.Size() > 1 {
		return base.Notify("Can't open for reading more than 1 chapter")
	}

	// use the toggled item, else the hovered one
	i := item
	if s.selected.Size() == 1 {
		i = s.selected.Keys()[0]
	}

	if i.readAvailablePath != "" {
		log.Log("Read format already downloaded")
		return s.readChapterCmd(ctx, i.readAvailablePath, i, config.ReadOptions())
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
	return s.downloadChapterCmd(ctx, i, downloadOptions, true)
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

func (s *state) anilistCmd() tea.Msg {
	mangaTitle := s.manga.Info().AnilistSearch
	if mangaTitle == "" {
		mangaTitle = s.manga.Info().Title
	}
	return anilistmangas.New(
		s.client.Anilist(),
		mangaTitle,
		func(manga lmanilist.Manga) tea.Cmd {
			return s.updateMetadataCmd(manga)
		},
	)
}
