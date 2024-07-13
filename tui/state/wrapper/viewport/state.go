package viewport

import (
	"context"
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

// State implements base.State. Wrapper of viewport.Model.
type State struct {
	viewport viewport.Model
	title    base.Title
	content  string

	borderHorizontalSize,
	borderVerticalSize int

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
	return s.title
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
func (s *State) Resize(size base.Size) tea.Cmd {
	s.viewport.Width = size.Width - s.borderHorizontalSize
	s.viewport.Height = size.Height - s.borderVerticalSize
	return nil
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return base.ShowLoadingMsg(false)
		},
		func() tea.Msg {
			return base.ShowSubtitleMsg(false)
		},
	)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.copy):
			return func() tea.Msg {
				return clipboard.WriteAll(s.content)
			}
		case key.Matches(msg, s.keyMap.goTop):
			s.viewport.GotoTop()
		case key.Matches(msg, s.keyMap.goBottom):
			s.viewport.GotoBottom()
		}
	case SetContentMsg:
		s.viewport.SetContent(string(msg))
		s.updateKeybinds()
	}
	s.viewport, cmd = s.viewport.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View() string {
	return s.viewport.View()
}

// updateKeybinds enables/disables keybinds based on the content.
func (s *State) updateKeybinds() {
	enable := s.content != ""
	s.keyMap.copy.SetEnabled(enable)
	s.keyMap.goTop.SetEnabled(enable)
	s.keyMap.goBottom.SetEnabled(enable)
}
