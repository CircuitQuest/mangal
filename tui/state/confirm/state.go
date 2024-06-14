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

var _ base.State = (*State)(nil)

type OnResponseFunc func(response bool) tea.Cmd

// State implements base.State.
type State struct {
	message    string
	keyMap     keyMap
	onResponse OnResponseFunc
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
	return base.Title{Text: "Confirm"}
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
		case key.Matches(msg, s.keyMap.yes):
			return s.onResponse(true)
		case key.Matches(msg, s.keyMap.no):
			return s.onResponse(false)
		}
	}

	return nil
}

// View implements base.State.
func (s *State) View() string {
	return fmt.Sprintf("%s %s", icon.Confirm, s.message)
}
