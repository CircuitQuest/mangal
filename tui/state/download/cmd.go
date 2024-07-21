package download

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata"
	"github.com/skratchdot/open-golang/open"
)

func (s *state) startDownloadCmd() tea.Msg {
	s.downloading = dSDownloading
	s.currentIdx = 0
	s.toDownload = s.chapters.ToDownload()
	return nextChapterMsg{}
}

func (s *state) downloadChapterCmd(ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		var (
			ch       = s.toDownload[s.currentIdx]
			downChap *metadata.DownloadedChapter
			err      error
		)

		downChap, err = s.client.DownloadChapter(ctx, ch.Chapter, s.options)
		if err != nil {
			errMsg := err.Error()
			// TODO: handle other responses here too if possible
			//
			// on too many requests error
			if strings.Contains(errMsg, "429") && strings.Contains(errMsg, "Retry-After") {
				s.retryCount++
				if s.retryCount > s.maxRetries {
					return fmt.Errorf("exceeded max retries (%d) while downloading chapters", s.maxRetries)
				}

				raTemp := strings.Split(errMsg, ":")
				raParsed, err := strconv.Atoi(strings.TrimSpace(raTemp[len(raTemp)-1]))
				if err != nil {
					return errors.New("error while parsing Retry-Count from error mesage: " + err.Error())
				}

				retryAfter := time.Duration(min(10, raParsed)) * time.Second
				return retryChapterMsg{After: retryAfter}
			}
		}
		ch.Down = downChap
		ch.Err = err

		if s.currentIdx+1 >= len(s.toDownload) {
			return downloadCompletedMsg{}
		}
		s.currentIdx++
		return nextChapterMsg{}
	}
}

// openCmd acts on the key press and thus relies on the keybinds being updated.
func (s *state) openCmd() tea.Msg {
	return open.Start(s.chapters[0].Down.Directory)
}

// retryCmd acts on the key press and thus relies on the keybinds being updated,
// will retry all failed chapters.
func (s *state) retryCmd() tea.Cmd {
	s.downloading = dSDownloading
	s.currentIdx = 0
	s.toDownload = s.chapters.Failed()
	return tea.Sequence(
		func() tea.Msg {
			return nextChapterMsg{}
		},
		s.Resize(s.size),
	)
}
