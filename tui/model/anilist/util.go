package anilist

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) updateCurrent() {
	switch m.current {
	case ID:
		m.input = &m.idInput
	case Secret:
		m.input = &m.secretInput
	case Code:
		m.input = &m.codeInput
	}
}

func (m *Model) viewLogin() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewField(ID),
		m.viewField(Secret),
		m.viewField(Code),
	)
}

func (m *Model) viewField(field Field) string {
	switch field {
	case ID:
		if m.current == field {
			return m.styles.selected.Render(m.selectCursor + m.idInput.View())
		}
		return m.styles.field.Render(m.idInput.View())
	case Secret:
		if m.current == field {
			return m.styles.selected.Render(m.selectCursor + m.secretInput.View())
		}
		return m.styles.field.Render(m.secretInput.View())
	case Code:
		if m.current == field {
			return m.styles.selected.Render(m.selectCursor + m.codeInput.View())
		}
		return m.styles.field.Render(m.codeInput.View())
	default:
		return ""
	}
}

func (m *Model) updateKeybinds() {
	m.keyMap.quit.SetEnabled(m.standalone && !m.inInput)

	m.keyMap.login.SetEnabled(!m.inInput && m.LoggedOut())
	m.keyMap.open.SetEnabled(!m.inInput && m.LoggedOut() && m.idInput.Value() != "")
	m.keyMap.logout.SetEnabled(!m.inInput && m.LoggedIn())

	navigate := m.LoggedOut() && !m.inInput
	m.keyMap.up.SetEnabled(navigate)
	m.keyMap.down.SetEnabled(navigate)
	m.keyMap.selekt.SetEnabled(navigate)

	m.keyMap.confirm.SetEnabled(m.inInput)
	m.keyMap.cancel.SetEnabled(m.inInput)
	m.keyMap.clear.SetEnabled(!m.inInput)
}

func sanitize(in string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(in)), "")
}
