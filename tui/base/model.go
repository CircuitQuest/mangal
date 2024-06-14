package base

import (
	"context"
	"time"

	"github.com/luevano/mangal/log"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zyedidia/generic/stack"
)

var _ tea.Model = (*model)(nil)

// model implements tea.model
//
// model is the parent of all States (windows), could be thought of as the main window.
type model struct {
	state   State
	history *stack.Stack[State]

	ctx       context.Context
	ctxCancel context.CancelFunc

	size Size

	styles styles

	keyMap *keyMap
	help   help.Model

	notification                string
	notificationDefaultDuration time.Duration

	// Custom states to show errors and logs
	errState func(err error) State
	logState func(title, content string, size Size) State
}

// Init implements tea.Model.
func (m *model) Init() tea.Cmd {
	return m.state.Init(m.ctx)
}

// stateSize returns the usable size of the state viewport (everything but the header/footer).
func (m *model) stateSize() Size {
	header := m.viewHeader()
	footer := m.viewFooter()

	size := m.size
	size.Height -= lipgloss.Height(header) + lipgloss.Height(footer)

	return size
}

func (m *model) cancel() {
	m.ctxCancel()
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
}

func (m *model) resize(size Size) {
	m.size = size
	m.help.Width = size.Width

	m.state.Resize(m.stateSize())
}

func (m *model) back() tea.Cmd {
	// do not pop the last state
	if m.history.Size() == 0 {
		return nil
	}

	log.L.Info().Str("state", m.history.Peek().Title().Text).Msg("going to the previous state")

	m.cancel()
	m.state = m.history.Pop()

	// update size for old models
	m.state.Resize(m.stateSize())

	return m.state.Init(m.ctx)
}

func (m *model) pushState(state State) tea.Cmd {
	log.L.Info().Str("state", state.Title().Text).Msg("new state")
	if !m.state.Intermediate() {
		m.history.Push(m.state)
	}

	m.state = state
	m.state.Resize(m.stateSize())

	return m.state.Init(m.ctx)
}
