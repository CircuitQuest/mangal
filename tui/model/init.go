package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init implements base.Model.
func (m *Model) Init() tea.Cmd {
	return m.state.Init(m)
}
