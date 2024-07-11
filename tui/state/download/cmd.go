package download

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
	"github.com/skratchdot/open-golang/open"
)

func (s *state) startDownloadCmd() tea.Msg {
	s.currentIdx = 0
	s.toDownload = s.chapters
	s.downloading = true
	return nextChapterMsg{}
}

func (s *state) downloadChapterCmd(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		chapter := s.toDownload[s.currentIdx]
		downChap, err := s.client.DownloadChapter(ctx, chapter, s.options)

		if err != nil {
			s.failed = append(s.failed, chapter)
		} else {
			s.succeedDownloads = append(s.succeedDownloads, downChap)
			s.succeed = append(s.succeed, chapter)
		}

		if s.currentIdx+1 >= len(s.toDownload) {
			return downloadCompletedMsg{}
		}

		s.currentIdx++
		return nextChapterMsg{}
	}
}

// openCmd relies on the keybinds enabled/disabled effectively
func (s *state) openCmd() tea.Msg {
	return open.Start(s.succeedDownloads[0].Directory)
}

// retryCmd relies on the keybinds enabled/disabled effectively
func (s *state) retryCmd() tea.Msg {
	s.currentIdx = 0
	s.toDownload = s.failed
	s.failed = []mangadata.Chapter{}
	s.downloading = true
	return nextChapterMsg{}
}
