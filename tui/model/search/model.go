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

// Unsearched is a convenience method to check if there hasn't been a search done.
func (m *Model) Unsearched() bool {
	return m.state == Unsearched
}

// Searching is a convenience method to check if a search is being performed.
func (m *Model) Searching() bool {
	return m.state == Searching
}

// Searched is a convenience method to check if a search was performed.
func (m *Model) Searched() bool {
	return m.state == Unsearched
}

// SearchCanceled is a convenience method to check if the search was canceled.
func (m *Model) SearchCanceled() bool {
	return m.state == SearchCanceled
}

// Focus sets the sate to Searching and enables the input.
func (m *Model) Focus() tea.Cmd {
	m.state = Searching
	m.input.CursorEnd()
	m.updateKeybinds()

	return tea.Sequence(
		m.input.Focus(),
		textinput.Blink,
	)
}

// Resize resizes the input given the new size.
func (m *Model) Resize(size base.Size) {
	m.input.Width = size.Width
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case base.RestoredMsg:
		return textinput.Blink
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.cancel):
			m.state = SearchCanceled
			m.input.Blur()
			m.input.Reset()
			m.updateKeybinds()
			// keep the last searched query in the input field
			m.input.SetValue(m.query)

			return SearchCancelCmd
		case key.Matches(msg, m.keyMap.confirm):
			if strings.TrimSpace(m.input.Value()) == "" {
				return base.Notify("Can't search whitespace only")
			}
			m.state = Searched
			m.input.Blur()
			m.updateKeybinds()

			m.query = m.input.Value()
			return SearchCmd(m.query)
		}
	}

	input, cmd := m.input.Update(msg)
	m.input = input
	return cmd
}

func (m *Model) View() string {
	return m.style.Render(m.input.View())
}

// updateKeybinds if the keymap should be enabled.
func (m *Model) updateKeybinds() {
	enable := m.state == Searching
	m.keyMap.cancel.SetEnabled(enable)
	m.keyMap.confirm.SetEnabled(enable)
}
