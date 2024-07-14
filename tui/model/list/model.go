package list

import (
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	stringutil "github.com/luevano/mangal/util/string"
)

type Model struct {
	list.Model
	delegate *list.DefaultDelegate
	keyMap   keyMap
}

func (m *Model) Backable() bool {
	return m.FilterState() == list.Unfiltered
}

func (m *Model) KeyMap() help.KeyMap {
	return m.keyMap
}

func (m *Model) Subtitle() string {
	singular, plural := m.StatusBarItemName()
	return stringutil.Quantify(len(m.VisibleItems()), singular, plural)
}

func (m *Model) Status() string {
	var p string
	if len(m.Items()) != 0 {
		p = m.Paginator.View()
	}

	if m.FilterState() == list.Filtering || m.FilterValue() != "" {
		return m.Paginator.View() + " " + m.FilterInput.View()
	}

	return p
}

func (m *Model) Resize(size base.Size) tea.Cmd {
	m.SetSize(size.Width, size.Height)
	return nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, m.keyMap.reverse):
			slices.Reverse(m.Items())
			return tea.Sequence(
				m.SetItems(m.Items()),
				base.Notify("Reversed"),
			)
		}
	}

end:
	list, cmd := m.Model.Update(msg)
	m.Model = list
	return cmd
}

func (m *Model) SetItems(items []list.Item) tea.Cmd {
	m.updateKeybinds(len(items) != 0)
	m.ResetSelected()
	return m.Model.SetItems(items)
}

// SetDelegateHeight sets the height of the delegate, which translates to the items' height.
//
// Clamps to a minimum of 1, in which case the description is hidden.
func (s *Model) SetDelegateHeight(height int) {
	if height < 2 {
		height = 1
	}
	if height == 1 {
		s.delegate.ShowDescription = false
	} else {
		s.delegate.ShowDescription = true
	}
	s.delegate.SetHeight(height)
	s.SetDelegate(s.delegate)
}

func (s *Model) updateKeybinds(enable bool) {
	s.keyMap.reverse.SetEnabled(enable)
}
