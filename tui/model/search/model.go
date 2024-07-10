package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/tui/base"
)

type State int

const (
	Unsearched State = iota
	Searching
	Searched
	SearchCanceled
)

var _ tea.Model = (*Model)(nil)

// Model implements tea.Model.
type Model struct {
	input textinput.Model
	state State
	query string

	style  lipgloss.Style
	keyMap keyMap
}

// Query returns the current set query (only when the input is confirmed)
func (m *Model) Query() string {
	return m.query
}

// State returns the current state of the search.
func (m *Model) State() State {
	return m.state
}

// Focus sets the sate to Searching and enables the input.
func (m *Model) Focus() tea.Cmd {
	m.state = Searching
	m.input.CursorEnd()
	m.enableKeyMap(true)

	return tea.Sequence(
		m.input.Focus(),
		textinput.Blink,
	)
}

// Resize resizes the input given the new size.
func (m *Model) Resize(size base.Size) {
	m.input.Width = size.Width
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case base.RestoredMsg:
		return m, textinput.Blink
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.cancel):
			m.state = SearchCanceled
			m.input.Blur()
			m.input.Reset()
			m.enableKeyMap(false)
			// keep the last searched query in the input field
			m.input.SetValue(m.query)

			return m, SearchCancelCmd
		case key.Matches(msg, m.keyMap.confirm):
			if strings.TrimSpace(m.input.Value()) == "" {
				return m, base.Notify("Can't search whitespace only")
			}
			m.state = Searched
			m.input.Blur()
			m.enableKeyMap(false)

			m.query = m.input.Value()
			return m, SearchCmd(m.query)
		}
	}

	input, updateCmd := m.input.Update(msg)
	searchChanged := m.input.Value() != input.Value()
	if searchChanged {
		m.keyMap.confirm.SetEnabled(m.input.Value() != "")
	}
	m.input = input

	return m, updateCmd
}

// View implements tea.Model.
func (m *Model) View() string {
	return m.style.Render(m.input.View())
}
