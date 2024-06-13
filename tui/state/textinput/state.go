package textinput

import (
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
func (s *State) Resize(size base.Size) {
	s.textinput.Width = size.Width
}

// Init implements base.State.
func (s *State) Init(model base.Model) tea.Cmd {
	return tea.Batch(s.textinput.Focus(), textinput.Blink)
}

// Update implements base.State.
func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.confirm) && strings.TrimSpace(s.textinput.Value()) != "":
			return s.options.OnResponse(strings.TrimSpace(s.textinput.Value()))
		}
	}

	s.textinput, cmd = s.textinput.Update(msg)
	return cmd
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.textinput.View()
}
