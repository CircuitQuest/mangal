package loading

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	title   string
	message string
	spinner spinner.Model
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
	return base.NoKeyMap{}
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: s.title, Background: color.Loading}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return ""
}

// Status implements base.State.
func (s *State) Status() string {
	return s.spinner.View()
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return s.spinner.Tick
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	s.spinner, cmd = s.spinner.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, style.Bold.Accent.Render(s.spinner.View()), style.Normal.Base.Render(s.message))
}

// SetMessage updates the message for the loading view.
func (s *State) SetMessage(message string) {
	s.message = message
}
