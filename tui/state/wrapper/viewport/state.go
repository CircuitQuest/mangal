package viewport

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

// State implements base.State. Wrapper of viewport.Model.
type State struct {
	viewport viewport.Model
	size     base.Size
	title    string
	content  string
	keyMap   keyMap
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
	return base.Title{Text: s.title, Background: color.Viewport}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return ""
}

// Status implements base.State.
func (s *State) Status() string {
	return fmt.Sprintf("%3.f%%", s.viewport.ScrollPercent()*100)
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.size = size
	s.viewport.Width = size.Width - 2 // -2 takes into account the border
	s.viewport.Height = size.Height - 2
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	s.viewport = viewport.New(s.size.Width-2, s.size.Height-2) // -2 takes into account the border
	s.viewport.SetContent(s.content)
	style := style.Normal.Base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(color.Viewport)
	s.viewport.Style = style
	s.keyMap = newKeyMap(s.viewport.KeyMap)
	return nil
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	s.viewport, cmd = s.viewport.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View() string {
	return s.viewport.View()
}
