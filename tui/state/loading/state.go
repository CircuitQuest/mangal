package loading

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	message  string
	subtitle string
	spinner  spinner.Model
	keyMap   help.KeyMap
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
	// return base.Title{Text: s.message, Background: color.Loading}
	return base.Title{Text: "Loading", Background: color.Loading}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return s.subtitle
}

// Status implements base.State.
func (s *State) Status() string {
	return s.spinner.View()
	// return ""
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
	return fmt.Sprint(
		style.Bold.Accent.Render(s.spinner.View()),
		style.Normal.Secondary.Render(s.message),
	)
}

// SetMessage updates the message for the loading view.
func (s *State) SetMessage(message string) {
	s.message = message
}
