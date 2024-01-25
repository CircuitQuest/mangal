package chapsdownloading

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/stringutil"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/ui/color"
	"github.com/luevano/mangal/ui/icon"
)

var _ base.State = (*State)(nil)

type State struct {
	options  Options
	chapters []libmangal.Chapter
	message  string

	progress progress.Model
	spinner  spinner.Model

	succeedPaths    []string
	succeed, failed []libmangal.Chapter
	currentIdx      int

	size base.Size

	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Downloading"}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
	s.size = size
	s.progress.Width = size.Width
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
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
				path, err := s.options.DownloadChapter(model.Context(), chapter)

				if err != nil {
					s.failed = append(s.failed, chapter)
				} else {
					s.succeedPaths = append(s.succeedPaths, path)
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
		return s.options.OnDownloadFinished(s.succeedPaths, s.succeed, s.failed)
	default:
		return nil
	}
}

func (s *State) View(model base.Model) string {
	spinnerView := s.spinner.View()
	return fmt.Sprintf(`%s Downloading %s - %d/%d

%s

%s %s`,
		icon.Progress,
		s.chapters[s.currentIdx].String(),
		s.currentIdx+1, len(s.chapters),
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.chapters))),
		spinnerView,
		lipgloss.NewStyle().Foreground(color.Secondary).Render(stringutil.Trim(s.message, s.size.Width-lipgloss.Width(spinnerView)-1)),
	)
}

func (s *State) Init(model base.Model) tea.Cmd {
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

func (s *State) SetMessage(message string) {
	s.message = message
}
