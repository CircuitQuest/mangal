package home

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/providers"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	providersState *providers.State
	keyMap         keyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return false
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "Mangal"}
}

// Subtitle implements base.State.
func (*State) Subtitle() string {
	return ""
}

// Status implements base.State.
func (s *State) Status() string {
	return ""
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) tea.Cmd {
	return nil
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	return tea.Sequence(
		base.Loading("Loading providers"),
		func() tea.Msg {
			loaders, err := manager.Loaders()
			if err != nil {
				return err
			}
			s.providersState = providers.New(loaders)
			return nil
		},
		base.Loaded,
		func() tea.Msg {
			if config.TUI.SkipHome.Get() {
				return s.providersState
			}
			return base.Notify("Providers loaded")()
		},
	)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.confirm):
			if s.providersState == nil {
				return base.Notify("Providers not yet loaded")
			}
			return func() tea.Msg {
				return s.providersState
			}
		}
	}
	return nil
}

// View implements base.State.
func (s *State) View() string {
	return meta.PrettyVersion()
}
