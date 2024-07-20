package anilist

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/util/cache"
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

func (m *Model) updateAuthUser() error {
	user, err := m.anilist.AuthenticatedUser()
	if err != nil {
		return err
	}
	m.user = user
	return nil
}

func (m *Model) updateUserHistory() error {
	err := cache.SetAnilistAuthHistory(m.userHistory)
	if err != nil {
		return err
	}
	items := make([]list.Item, m.userHistory.Size())
	for i, u := range m.userHistory {
		items[i] = &item{user: u}
	}
	m.list.SetItems(items)
	return nil
}

func (m *Model) viewLoggedIn() string {
	if m.inNew {
		return m.viewLoggedOut()
	}
	view := "Logged in to " + m.user.Name + " (" + strconv.Itoa(m.user.ID) + ")"
	if !m.standalone {
		return view
	}

	msg := m.styles.notification.Render(m.notification)
	view = lipgloss.JoinVertical(lipgloss.Left, m.title, msg, view, m.help.View(&m.keyMap))
	return m.styles.view.Render(view)
}

func (m *Model) viewLoggedOut() string {
	var k help.KeyMap
	k = base.CombinedKeyMap(&m.keyMap, &m.list.KeyMap)
	view := m.viewAvailableLogins()
	if m.inNew {
		k = &m.keyMap
		view = m.viewNewLogin()
	}
	if !m.standalone {
		return view
	}

	msg := m.styles.notification.Render(m.notification)
	return lipgloss.JoinVertical(lipgloss.Left, m.title, msg, view, m.help.View(k))
}

// viewAvailableLogins renders the list of available logins.
func (m *Model) viewAvailableLogins() string {
	return lipgloss.JoinVertical(lipgloss.Left, "Available logins", m.list.View())
}

// viewNewLogin renders the login form.
func (m *Model) viewNewLogin() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewField(ID),
		m.viewField(Secret),
		m.viewField(Code),
	)
}

// viewField renders the input fields depending and render hovered accordingly.
func (m *Model) viewField(field Field) string {
	cursor := ""
	render := m.styles.field.Render
	if m.current == field {
		cursor = m.selectCursor
		render = m.styles.selected.Render
	}

	switch field {
	case ID:
		return render(cursor + m.idInput.View())
	case Secret:
		return render(cursor + m.secretInput.View())
	case Code:
		return render(cursor + m.codeInput.View())
	default:
		return "unkown field"
	}
}

// updateKeybinds enables/disables the keybinds depending on the state of
// the anilist authentication status.
func (m *Model) updateKeybinds() {
	m.keyMap.quit.SetEnabled(m.standalone && !m.inInput)
	m.keyMap.back.SetEnabled(m.inNew && !m.inInput && m.userHistory.Size() != 0)

	loggable := m.LoggedOut() && !m.inInput
	m.keyMap.login.SetEnabled(loggable)
	m.keyMap.open.SetEnabled(loggable && m.inNew && m.idInput.Value() != "")
	m.keyMap.logout.SetEnabled(m.LoggedIn())
	m.keyMap.new.SetEnabled(m.LoggedOut() && !m.inNew)
	m.keyMap.delete.SetEnabled(m.LoggedOut() && !m.inNew && m.userHistory.Size() != 0)

	navigate := m.inNew && !m.inInput
	m.keyMap.up.SetEnabled(navigate)
	m.keyMap.down.SetEnabled(navigate)
	m.keyMap.selekt.SetEnabled(navigate)
	m.keyMap.clear.SetEnabled(navigate)

	typing := m.inNew && m.inInput
	m.keyMap.confirm.SetEnabled(typing)
	m.keyMap.cancel.SetEnabled(typing)
}

func sanitize(in string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(in)), "")
}
