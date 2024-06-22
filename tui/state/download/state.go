package download

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/skratchdot/open-golang/open"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	client   *libmangal.Client
	chapters []mangadata.Chapter
	options  libmangal.DownloadOptions

	// currently downloading?
	downloading bool

	currentIdx int
	toDownload,
	succeed,
	failed []mangadata.Chapter
	succeedDownloads []*metadata.DownloadedChapter

	message  string
	progress progress.Model
	spinner  spinner.Model

	size   base.Size
	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return true
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return true
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.keyMap
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{Text: "Download"}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	return ""
}

// Status implements base.State.
func (s *state) Status() string {
	return ""
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	s.size = size
	s.progress.Width = size.Width
	return nil
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	s.client.Logger().SetOnLog(func(format string, a ...any) {
		s.message = fmt.Sprintf(format, a...)
		log.Log(format, a...)
	})

	s.setKeyMapEnabled(false)
	return tea.Sequence(
		s.spinner.Tick,
		startDownloadCmd,
	)
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.open) && len(s.succeedDownloads) > 0:
			err := open.Start(s.succeedDownloads[0].Directory)
			// err := open.Start(filepath.Dir(s.options.SucceedDownloads[0].Path()))
			if err != nil {
				return func() tea.Msg {
					return err
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.retry) && len(s.failed) > 0:
			return downloadFailedCmd
		}
	case spinner.TickMsg:
		spinner, cmd := s.spinner.Update(msg)
		s.spinner = spinner
		return cmd
	case startDownloadMsg:
		s.currentIdx = 0
		s.downloading = true
		s.toDownload = s.chapters
		return nextChapterCmd
	case downloadFailedMsg:
		s.currentIdx = 0
		s.downloading = true
		s.toDownload = s.failed
		s.failed = []mangadata.Chapter{}
		return nextChapterCmd
	case nextChapterMsg:
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
	case downloadCompletedMsg:
		s.downloading = false
		s.setKeyMapEnabled(true)
	}
	return nil
}

// View implements base.State.
func (s *state) View() string {
	if s.downloading {
		return s.viewDownloading()
	}

	return s.viewDownloaded()
}

func (s *state) setKeyMapEnabled(enabled bool) {
	s.keyMap.open.SetEnabled(enabled)
	// retry should only be enabled if there are failed downloads
	s.keyMap.retry.SetEnabled(len(s.failed) > 0)
}

func (s *state) viewDownloading() string {
	spinnerView := s.spinner.View()
	return fmt.Sprintf(`%s Downloading %s - %d/%d

%s

%s %s`,
		icon.Progress,
		s.toDownload[s.currentIdx].String(),
		s.currentIdx+1, len(s.toDownload),
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.toDownload))),
		spinnerView,
		style.Normal.Secondary.Render(stringutil.Trim(s.message, s.size.Width-lipgloss.Width(spinnerView)-1)),
	)
}

func (s *state) viewDownloaded() string {
	var (
		succeed = len(s.succeed)
		failed  = len(s.failed)
	)

	if failed == 0 {
		return style.Normal.Success.
			Render(fmt.Sprintf(
				"%s downloaded successfully!",
				stringutil.Quantify(succeed, "chapter", "chapters"),
			))
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"%s downloaded successfully, %d failed.",
		stringutil.Quantify(succeed, "chapter", "chapters"),
		failed,
	))

	sb.WriteString("\n\nFailed:\n")

	if failed <= 3 {
		for _, chapter := range s.failed {
			sb.WriteString(fmt.Sprintf("\n%s", chapter))
		}
	} else {
		indices := make([]float32, failed)
		for i, c := range s.failed {
			indices[i] = c.Info().Number
		}

		sb.WriteString(stringutil.FormatRanges(indices))
	}

	return sb.String()
}
