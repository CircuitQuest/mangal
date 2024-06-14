package chapsdownloaded

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/skratchdot/open-golang/open"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	options Options
	keyMap  keyMap
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
	return base.Title{Text: "Done"}
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
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return nil
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.quit):
			return tea.Quit
		case key.Matches(msg, s.keyMap.open) && len(s.options.SucceedDownloads) > 0:
			err := open.Start(s.options.SucceedDownloads[0].Directory)
			// err := open.Start(filepath.Dir(s.options.SucceedDownloads[0].Path()))
			if err != nil {
				return func() tea.Msg {
					return err
				}
			}

			return nil
		case key.Matches(msg, s.keyMap.retry) && len(s.options.Failed) > 0:
			return s.options.DownloadChapters(s.options.Failed)
		}
	}

	return nil
}

// View implements base.State.
func (s *State) View() string {
	var (
		succeed = len(s.options.Succeed)
		failed  = len(s.options.Failed)
	)

	if len(s.options.Failed) == 0 {
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
		for _, chapter := range s.options.Failed {
			sb.WriteString(fmt.Sprintf("\n%s", chapter))
		}
	} else {
		indices := make([]float32, failed)
		for i, c := range s.options.Failed {
			indices[i] = c.Info().Number
		}

		sb.WriteString(stringutil.FormatRanges(indices))
	}

	return sb.String()
}
