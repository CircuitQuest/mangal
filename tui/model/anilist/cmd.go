package anilist

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/skratchdot/open-golang/open"
)

func (m *Model) loginCachedCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		ok, err := m.anilist.AuthorizeCachedUser(item.user)
		if err != nil {
			return err
		}
		if !ok {
			m.userHistory.Delete(item.user)
			err := m.updateUserHistory()
			if err != nil {
				return err
			}
			return NotificationMsg(fmt.Sprintf("cached user %q not found, deleted from history", item.user))
		}

		// Update auth user
		err = m.updateAuthUser()
		if err != nil {
			return err
		}

		// Update user history
		m.userHistory.Add(m.user.Name)
		err = m.updateUserHistory()
		if err != nil {
			return err
		}

		m.state = LoggedIn
		m.inNew = false
		m.updateKeybinds()
		return NotificationMsg(fmt.Sprintf("cached user %q logged in to anilist", item.user))
	}
}

func (m *Model) loginCmd() tea.Msg {
	id := sanitize(m.idInput.Value())
	secret := sanitize(m.secretInput.Value())
	code := sanitize(m.codeInput.Value())
	if id == "" {
		return NotificationMsg("anilist login error: ID is empty")
	}
	if code == "" {
		return NotificationMsg("anilist login error: Code is empty")
	}

	var err error
	// Secret empty means an implicit grant (directly with access token)
	if secret == "" {
		err = m.anilist.AuthorizeWithAccessToken(context.Background(), code)
	} else {
		err = m.anilist.AuthorizeWithCodeGrant(context.Background(), anilist.CodeGrant{
			ID:     id,
			Secret: secret,
			Code:   code,
		})
	}

	if err != nil {
		return NotificationMsg("anilist login error: " + err.Error())
	}

	// Update auth user
	err = m.updateAuthUser()
	if err != nil {
		return nil
	}

	// Update user history
	m.userHistory.Add(m.user.Name)
	err = m.updateUserHistory()
	if err != nil {
		return err
	}

	m.state = LoggedIn
	m.inNew = false
	m.updateKeybinds()
	return NotificationMsg(fmt.Sprintf("user %q logged in to anilist", m.user.Name))
}

func (m *Model) logoutCmd() tea.Msg {
	if err := m.anilist.Logout(false); err != nil {
		return NotificationMsg("anilist logout error: " + err.Error())
	}

	username := m.user.Name
	m.user = anilist.User{}

	m.state = LoggedOut
	m.updateKeybinds()
	return NotificationMsg(fmt.Sprintf("user %q logged out of anilist", username))
}

func (m *Model) deleteUserCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		// Remove auth user data
		err := m.anilist.DeleteCachedUser(item.user)
		if err != nil {
			return err
		}

		// Update user history
		m.userHistory.Delete(item.user)
		err = m.updateUserHistory()
		if err != nil {
			return err
		}
		if m.userHistory.Size() == 0 {
			m.inNew = true
		}
		m.updateKeybinds()
		return NotificationMsg(fmt.Sprintf("deleted user %q from cache", item.user))
	}
}

// updateInputCmd sets the current input command to the one hovered
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
	codeType := "code"
	if m.secretInput.Value() == "" {
		codeType = "token"
	}
	err := open.Run(fmt.Sprintf(authURL, m.idInput.Value(), codeType))
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
