package errorstate

import (
	"context"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	error  error
	size   base.Size
	keyMap keyMap
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
	return base.Title{Text: "Error", Background: color.Error}
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
		case key.Matches(msg, s.keyMap.copyError):
			return func() tea.Msg {
				return clipboard.WriteAll(s.error.Error())
			}
		}
	}

	return nil
}

// View implements base.State.
func (s *State) View() string {
	return style.Normal.Error.Width(s.size.Width).Height(s.size.Height).Render(s.error.Error())
}
