package base

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/log"
	"github.com/pkg/errors"
)

// Update implements tea.Model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, m.resize(Size{Width: msg.Width, Height: msg.Height})
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.back) && m.state.Backable():
			return m, m.back()
		case key.Matches(msg, m.keyMap.help):
			return m, m.toggleHelp()
		case key.Matches(msg, m.keyMap.log):
			return m, m.pushState(m.logState("Logs", log.Aggregate.String()))
		}
	case NotificationMsg:
		return m, m.notify(msg.Message, m.notificationDefaultDuration)
	case NotificationWithDurationMsg:
		return m, m.notify(msg.Message, msg.Duration)
	case NotificationTimeoutMsg:
		return m, m.hideNotification()
	case BackMsg:
		return m, m.back()
	case State:
		return m, m.pushState(msg)
	case error:
		if errors.Is(msg, context.Canceled) || strings.Contains(msg.Error(), context.Canceled.Error()) {
			return m, nil
		}

		log.L.Err(msg).Msg("")

		return m, m.pushState(m.errState(msg))
	}

	cmd := m.state.Update(m.ctx, msg)
	return m, cmd
}

// resize the program
func (m *model) resize(size Size) tea.Cmd {
	m.size = size
	m.help.Width = size.Width

	return m.resizeState()
}

// resizeState only resizes the State with its updated size
func (m *model) resizeState() tea.Cmd {
	return m.state.Resize(m.stateSize())
}

// back to the previous State
func (m *model) back() tea.Cmd {
	// do not pop the last state
	if m.history.Size() == 0 {
		return nil
	}

	log.L.Info().Str("state", m.history.Peek().Title().Text).Msg("going to the previous state")

	// TODO: perform m.state.Destroy() once implemented?
	m.cancel()
	m.state = m.history.Pop()

	return m.resizeState()
}

// pushState initializes a new state and pushes previous one into history
func (m *model) pushState(state State) tea.Cmd {
	log.L.Info().Str("state", state.Title().Text).Msg("new state")
	if !m.state.Intermediate() {
		m.history.Push(m.state)
	}
	m.state = state

	return tea.Sequence(
		m.resizeState(),
		m.state.Init(m.ctx),
	)
}

// toggleHelp show/hide help menu (keybindings)
func (m *model) toggleHelp() tea.Cmd {
	m.help.ShowAll = !m.help.ShowAll
	return m.resizeState()
}

// notify will show a message to the right of the title and status.
func (m *model) notify(message string, duration time.Duration) tea.Cmd {
	m.notification = message
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return NotificationTimeoutMsg{}
	})
}

// hideNotification removes the notification
func (m *model) hideNotification() tea.Cmd {
	m.notification = ""
	return nil
}
