package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

	maxWidth int
	// the Height field in size is used
	// for the max suggestions to display
	size   base.Size
	styles styles
	keyMap keyMap
}

// Query returns the current set query (only when the input is confirmed)
func (m *Model) Query() string {
	return m.query
}

// SetSuggestions sets the suggestions for the internal input.
func (m *Model) SetSuggestions(suggestions []string) {
	m.input.SetSuggestions(suggestions)
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

	return m.input.Focus()
}

// Resize resizes the input given the new size.
func (m *Model) Resize(size base.Size) {
	if size.Width < m.maxWidth {
		m.size.Width = size.Width
		m.input.Width = size.Width
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.confirm):
			// Remove all surrounding whitespace
			q := strings.TrimSpace(m.input.Value())
			if q == "" {
				return base.Notify("Can't search whitespace only")
			}

			m.state = Searched
			m.input.Blur()
			m.updateKeybinds()

			// remove redundant spaces and use that value as the query
			q = strings.Join(strings.Fields(q), " ")
			m.input.SetValue(q)
			m.query = q
			return SearchCmd(q)
		case key.Matches(msg, m.keyMap.cancel):
			m.state = SearchCanceled
			m.input.Blur()
			m.input.Reset()
			m.updateKeybinds()
			// keep the last searched query in the input field
			m.input.SetValue(m.query)

			return SearchCancelCmd
		}
	}

	input, cmd := m.input.Update(msg)
	m.input = input
	return cmd
}

func (m *Model) View() string {
	return m.input.View()
}

// SuggestionBox returns the rendered suggestions box with styling applied.
func (m *Model) SuggestionBox() string {
	sugs, idx := m.getSuggestions()
	if len(sugs) == 0 {
		return m.styles.renderSuggestionBox(m.styles.normalSuggestion.Render("<no suggestions>"), m.size.Width)
	}
	if len(sugs) > m.size.Height {
		sugs = sugs[:m.size.Height]
	}

	var sb strings.Builder
	sb.Grow(200)
	for i, s := range sugs {
		if i == idx {
			sb.WriteString(m.styles.matchSuggestion.Render(s))
		} else {
			sb.WriteString(m.styles.normalSuggestion.Render(s))
		}
		if i < len(sugs)-1 {
			sb.WriteByte('\n')
		}
	}
	return m.styles.renderSuggestionBox(sb.String(), m.size.Width)
}

// TODO: change to exposed input methods once (if?) merged and released,
// https://github.com/charmbracelet/bubbles/pull/556
//
// Also, CurrentSuggestion is currently panic if called on empty
// matches (no input or no matches at all); fixed by (not yet released)
// https://github.com/charmbracelet/bubbles/pull/473
//
// getSuggestions will get the matched suggestions and their index, else
// just the available suggestions and -1 as index
func (m *Model) getSuggestions() (matches []string, idx int) {
	idx = -1
	val := m.input.Value()
	if strings.TrimSpace(val) == "" {
		return m.input.AvailableSuggestions(), idx
	}

	i := 0
	for _, s := range m.input.AvailableSuggestions() {
		if strings.HasPrefix(strings.ToLower(s), strings.ToLower(val)) {
			if strings.ToLower(s) == strings.ToLower(m.input.CurrentSuggestion()) {
				idx = i
			}
			matches = append(matches, s)
			i++
		}
	}
	return matches, idx
}

// updateKeybinds if the keymap should be enabled.
func (m *Model) updateKeybinds() {
	enable := m.Searching()
	m.keyMap.cancel.SetEnabled(enable)
	m.keyMap.confirm.SetEnabled(enable)
}
