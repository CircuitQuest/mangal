package model

import (
	"context"

	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zyedidia/generic/stack"
)

var (
	_ base.Model  = (*Model)(nil)
	_ help.KeyMap = (*Model)(nil)
)

// Model implements base.Model
//
// Model is the parent of all States (windows), could be thought of as the main window.
type Model struct {
	state   base.State
	history *stack.Stack[base.State]

	context           context.Context
	contextCancelFunc context.CancelFunc

	size base.Size

	styles base.Styles

	keyMap *keyMap
	help   help.Model
}

// Context implements base.Model.
func (m *Model) Context() context.Context {
	return m.context
}

// ShortHelp implements help.KeyMap.
func (m *Model) ShortHelp() []key.Binding {
	keys := []key.Binding{m.keyMap.back, m.keyMap.help}
	return append(keys, m.state.KeyMap().ShortHelp()...)
}

// FullHelp implements help.KeyMap.
func (m *Model) FullHelp() [][]key.Binding {
	keys := [][]key.Binding{{m.keyMap.back, m.keyMap.help, m.keyMap.quit, m.keyMap.log}}
	return append(keys, m.state.KeyMap().FullHelp()...)
}

func (m *Model) StateSize() base.Size {
	var height int

	if m.help.ShowAll {
		height = m.size.Height - lipgloss.Height(m.help.View(m)) - 2
	} else {
		height = m.size.Height - 3
	}

	if m.state.Subtitle() != "" {
		height -= 2
	}

	return base.Size{
		Width:  m.size.Width,
		Height: height,
	}
}

func (m *Model) cancel() {
	m.contextCancelFunc()
	m.context, m.contextCancelFunc = context.WithCancel(context.Background())
}

func (m *Model) resize(size base.Size) {
	m.size = size
	m.help.Width = size.Width

	m.state.Resize(m.StateSize())
}

func (m *Model) back() tea.Cmd {
	// do not pop the last state
	if m.history.Size() == 0 {
		return nil
	}

	log.L.Info().Str("state", m.history.Peek().Title().Text).Msg("going to the previous state")

	m.cancel()
	m.state = m.history.Pop()

	// update size for old models
	m.state.Resize(m.StateSize())

	return m.state.Init(m)
}

func (m *Model) pushState(state base.State) tea.Cmd {
	log.L.Info().Str("state", state.Title().Text).Msg("new state")
	if !m.state.Intermediate() {
		m.history.Push(m.state)
	}

	m.state = state
	m.state.Resize(m.StateSize())

	return m.state.Init(m)
}
