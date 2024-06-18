package textinput

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*State)(nil)

type OnResponseFunc func(response string) tea.Cmd

// State implements base.State. Wrapper of textinput.Model.
type State struct {
	options Options

	textinput textinput.Model
	keyMap    keyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return s.options.Intermediate
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
	return s.options.Title
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return s.options.Subtitle
}

// Status implements base.State.
func (s *State) Status() string {
	return ""
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) tea.Cmd {
	s.textinput.Width = size.Width
	return nil
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return tea.Batch(s.textinput.Focus(), textinput.Blink)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.confirm) && strings.TrimSpace(s.textinput.Value()) != "":
			s.textinput.Blur()
			return tea.Sequence(
				s.options.OnResponse(strings.TrimSpace(s.textinput.Value())),
				s.Init(ctx), // re-enable the prompt after the response, so that it's usable when backing up
			)
		}
	}

	s.textinput, cmd = s.textinput.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View() string {
	return s.textinput.View()
}
