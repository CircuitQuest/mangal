package download

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
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

	s.updateKeyMap()
	return tea.Sequence(
		s.spinner.Tick,
		s.startDownloadCmd,
	)
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.open):
			return s.openCmd
		case key.Matches(msg, s.keyMap.retry):
			return s.retryCmd
		}
	case spinner.TickMsg:
		spinner, cmd := s.spinner.Update(msg)
		s.spinner = spinner
		return cmd
	case nextChapterMsg:
		return s.downloadChapterCmd(ctx)
	case downloadCompletedMsg:
		s.downloading = false
		s.updateKeyMap()
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
