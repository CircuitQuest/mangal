package download

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skratchdot/open-golang/open"
)

func (s *state) startDownloadCmd() tea.Msg {
	s.downloading = dSDownloading
	s.currentIdx = 0
	s.toDownload = s.chapters.toDownload()
	return nextChapterMsg{}
}

func (s *state) downloadChapterCmd(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		ch := s.toDownload[s.currentIdx]
		downChap, err := s.client.DownloadChapter(ctx, ch.chapter, s.options)
		ch.down = downChap

		if err != nil {
			ch.state = cSFailed
		} else {
			ch.state = cSSucceed
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
	return open.Start(s.chapters[0].down.Directory)
}

// retryCmd relies on the keybinds enabled/disabled effectively
func (s *state) retryCmd() tea.Cmd {
	s.downloading = dSDownloading
	s.currentIdx = 0
	s.toDownload = s.chapters.failed()
	return tea.Sequence(
		func() tea.Msg {
			return nextChapterMsg{}
		},
		s.Resize(s.size),
	)
}
