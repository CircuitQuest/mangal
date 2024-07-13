package confirm

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/tui/base"
)

var _ base.State = (*state)(nil)

type onResponseFunc func(response bool) tea.Cmd

// state implements base.state.
type state struct {
	message    string
	keyMap     keyMap
	onResponse onResponseFunc
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
	return base.Title{Text: "Confirm"}
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
	return nil
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	return nil
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.yes):
			return s.onResponse(true)
		case key.Matches(msg, s.keyMap.no):
			return s.onResponse(false)
		}
	}

	return nil
}

// View implements base.State.
func (s *state) View() string {
	return fmt.Sprintf("%s %s", icon.Confirm.Colored(), s.message)
}
