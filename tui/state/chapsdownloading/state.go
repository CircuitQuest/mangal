package chapsdownloading

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	stringutil "github.com/luevano/mangal/util/string"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	options  Options
	chapters []mangadata.Chapter
	message  string

	progress progress.Model
	spinner  spinner.Model

	succeedDownloads []*metadata.DownloadedChapter
	succeed, failed  []mangadata.Chapter
	currentIdx       int

	size base.Size

	keyMap help.KeyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return true
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return true
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "Downloading"}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return ""
}

// Status implements base.State.
func (s *State) Status() string {
	return ""
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.size = size
	s.progress.Width = size.Width
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return nextChapterIdxMsg(0)
		},
		func() tea.Msg {
			return spinner.TickMsg{}
		},
		s.progress.Init(),
	)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case progress.FrameMsg:
		progressModel, cmd := s.progress.Update(msg)
		s.progress = progressModel.(progress.Model)
		return cmd
	case spinner.TickMsg:
		s.spinner, cmd = s.spinner.Update(msg)
		return cmd
	case nextChapterIdxMsg:
		s.currentIdx = int(msg)
		return tea.Sequence(
			func() tea.Msg {
				chapter := s.chapters[msg]
				downChap, err := s.options.DownloadChapter(ctx, chapter)

				if err != nil {
					s.failed = append(s.failed, chapter)
				} else {
					s.succeedDownloads = append(s.succeedDownloads, downChap)
					s.succeed = append(s.succeed, chapter)
				}

				nextIndex := msg + 1

				if int(nextIndex) >= len(s.chapters) {
					return downloadCompletedMsg{}
				}

				return nextIndex
			},
		)
	case downloadCompletedMsg:
		return s.options.OnDownloadFinished(s.succeedDownloads, s.succeed, s.failed)
	default:
		return nil
	}
}

// View implements base.State.
func (s *State) View() string {
	spinnerView := s.spinner.View()
	return fmt.Sprintf(`%s Downloading %s - %d/%d

%s

%s %s`,
		icon.Progress,
		s.chapters[s.currentIdx].String(),
		s.currentIdx+1, len(s.chapters),
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.chapters))),
		spinnerView,
		style.Normal.Secondary.Render(stringutil.Trim(s.message, s.size.Width-lipgloss.Width(spinnerView)-1)),
	)
}

// SetMessage updates the message for the loading view.
func (s *State) SetMessage(message string) {
	s.message = message
}
