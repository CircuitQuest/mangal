package metadata

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/metadata"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	viewport base.State
	meta     *metadata.Model

	enumeratorStyle lipgloss.Style
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
	return base.Title{
		Text:       s.meta.Style().Prefix + " Metadata",
		Background: s.meta.Style().Color,
		Foreground: color.Bright,
	}
}

// Subtitle implements base.state.
func (s *State) Subtitle() string {
	return ""
}

// Status implements base.state.
func (s *State) Status() string {
	return s.meta.View() + " " + s.viewport.Status()
}

// Resize implements base.state.
func (s *State) Resize(size base.Size) tea.Cmd {
	return s.viewport.Resize(size)
}

// Init implements base.state.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return tea.Sequence(
		s.viewport.Init(ctx),
		func() tea.Msg {
			return viewport.SetContentMsg(s.renderMetadata())
		},
	)
}

// Update implements base.state.
func (s *State) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	return s.viewport.Update(ctx, msg)
}

// View implements base.state.
func (s *State) View() string {
	return s.viewport.View()
}
