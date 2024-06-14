package base

import (
	"context"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/log"
	"github.com/pkg/errors"
)

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(Size{
			Width:  msg.Width,
			Height: msg.Height,
		})

		return m, nil
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.back) && m.state.Backable():
			return m, m.back()
		case key.Matches(msg, m.keyMap.help):
			m.help.ShowAll = !m.help.ShowAll
			m.resize(m.size)
			return m, nil
		case key.Matches(msg, m.keyMap.log):
			return m, m.pushState(m.logState("Logs", log.Aggregate.String(), m.stateSize()))
		}
	case BackMsg:
		// this msg can override Backable() output
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
