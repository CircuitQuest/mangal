package anilist

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/skratchdot/open-golang/open"
)

func (m *Model) loginCmd() tea.Msg {
	id := sanitize(m.idInput.Value())
	secret := sanitize(m.secretInput.Value())
	code := sanitize(m.codeInput.Value())
	if id == "" {
		return NotificationMsg("anilist login error: ID is empty")
	}
	if secret == "" {
		return NotificationMsg("anilist login error: Secret is empty")
	}
	if code == "" {
		return NotificationMsg("anilist login error: Code is empty")
	}
	err := m.anilist.Authorize(context.Background(), anilist.LoginCredentials{
		ID:     id,
		Secret: secret,
		Code:   code,
	})
	if err != nil {
		return NotificationMsg("anilist login error: " + err.Error())
	}

	m.state = LoggedIn
	m.updateKeybinds()
	return NotificationMsg("successfully logged in to anilist")
}

func (m *Model) logoutCmd() tea.Msg {
	if err := m.anilist.Logout(); err != nil {
		return NotificationMsg("anilist logout error: " + err.Error())
	}
	m.state = LoggedOut
	m.updateKeybinds()
	return NotificationMsg("successfully logged out of anilist")
}

func (m *Model) updateInputCmd(msg tea.Msg) tea.Cmd {
	input, cmd := m.input.Update(msg)
	switch m.current {
	case ID:
		m.idInput = input
	case Secret:
		m.secretInput = input
	case Code:
		m.codeInput = input
	}
	m.updateCurrent()
	return cmd
}

func (m *Model) openAuthURLCmd() tea.Msg {
	err := open.Run(fmt.Sprintf(authURL, m.idInput.Value()))
	if err != nil {
		return err
	}
	return NotificationMsg("opening auth url")
}

func (m *Model) notifyCmd(message string) tea.Cmd {
	m.notification = message
	return tea.Tick(m.notificationDuration, func(t time.Time) tea.Msg {
		return NotificationTimeoutMsg{}
	})
}
