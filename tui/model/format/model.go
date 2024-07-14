package format

import (
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/list"
)

type forWhat string

const (
	forRead     forWhat = "read"
	forDownload forWhat = "download"
	forBoth     forWhat = "both"
)

var _ tea.Model = (*Model)(nil)

// Model implements tea.Model.
type Model struct {
	list *list.Model

	// pre-rendered
	title,
	help string

	minSize,
	size,
	maxSize base.Size

	keyMap *keyMap
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return m.list.Resize(m.size)
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == _list.Filtering {
			goto end
		}

		// guaranteed to always have items
		i, _ := m.list.SelectedItem().(*item)
		switch {
		case key.Matches(msg, m.keyMap.setRead):
			return m, m.setFormatForCmd(forRead, i.format)
		case key.Matches(msg, m.keyMap.setDownload):
			return m, m.setFormatForCmd(forDownload, i.format)
		case key.Matches(msg, m.keyMap.setBoth):
			return m, m.setFormatForCmd(forBoth, i.format)
		case key.Matches(msg, m.keyMap.back):
			return m, backCmd
		}
	}
end:
	// the ctx is not used for anything in the list anyways
	return m, m.list.Update(msg)
}

// View implements tea.Model.
func (m *Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title, m.list.View(), m.help)
}
