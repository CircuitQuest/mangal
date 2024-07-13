package metadata

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	viewport *viewport.State
	meta     metadata.Metadata

	styles styles
}

// Intermediate implements base.state.
func (s *State) Intermediate() bool {
	return true
}

// Backable implements base.state.
func (s *State) Backable() bool {
	return true
}

// KeyMap implements base.state.
func (s *State) KeyMap() help.KeyMap {
	return s.viewport.KeyMap()
}

// Title implements base.state.
func (s *State) Title() base.Title {
	return s.viewport.Title()
}

// Subtitle implements base.state.
func (s *State) Subtitle() string {
	return ""
}

// Status implements base.state.
func (s *State) Status() string {
	return s.viewport.Status()
}

// Resize implements base.state.
func (s *State) Resize(size base.Size) tea.Cmd {
	return s.viewport.Resize(size)
}

// Init implements base.state.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return s.viewport.Init(ctx)
}

// Update implements base.state.
func (s *State) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	return s.viewport.Update(ctx, msg)
}

// View implements base.state.
func (s *State) View() string {
	return s.viewport.View()
}
