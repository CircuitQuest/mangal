package anilist

import (
	"errors"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/util/cache"
)

var authURL = "https://anilist.co/api/v2/oauth/authorize?client_id=%s&response_type=%s"

type State int

const (
	Uninitialized State = iota
	LoggedIn
	LoggedOut
)

type Field int

const (
	ID Field = iota
	Secret
	Code
)

var _ tea.Model = (*Model)(nil)

type Model struct {
	idInput,
	secretInput,
	codeInput textinput.Model
	input   *textinput.Model // current selected input
	list    *list.Model
	help    help.Model
	anilist *anilist.Anilist

	user        anilist.User
	userHistory cache.UserHistory

	// to be able to clear the screen in standalone
	quitting,
	standalone bool

	notification         string
	notificationDuration time.Duration

	// pre-rendered
	title,
	selectCursor string

	state   State
	current Field
	inInput,
	inNew bool

	styles styles
	keyMap keyMap
}

// State returns the current state.
func (m *Model) State() State {
	return m.state
}

// Uninitialized is a convenience method to check if the model hasn't been initialized..
func (m *Model) Uninitialized() bool {
	return m.state == Uninitialized
}

// LoggedIn is a convenience method to check if logged into anilist.
func (m *Model) LoggedIn() bool {
	return m.state == LoggedIn
}

// LoggedOut is a convenience method to check if logged out of anilist.
func (m *Model) LoggedOut() bool {
	return m.state == LoggedOut
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	if m.anilist == nil {
		m.notification = "anilist login: anilist is nil"
		return nil
	}
	if m.anilist.Authenticated() {
		m.state = LoggedIn
	} else {
		m.state = LoggedOut
	}
	if m.userHistory.Size() == 0 {
		m.inNew = true
	}
	m.updateKeybinds()
	return tea.Sequence(
		m.list.Init(),
		textinput.Blink,
	)
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// selected user from history
		i, _ := m.list.SelectedItem().(*item)
		switch {
		case key.Matches(msg, m.keyMap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.back):
			m.inNew = false
			m.updateKeybinds()
		case key.Matches(msg, m.keyMap.up):
			m.current--
			if m.current < ID {
				m.current = Code
			}
			m.updateCurrent()
		case key.Matches(msg, m.keyMap.down):
			m.current++
			if m.current > Code {
				m.current = ID
			}
			m.updateCurrent()
		case key.Matches(msg, m.keyMap.selekt):
			m.inInput = true
			m.input.CursorEnd()
			m.updateKeybinds()
			return m, m.input.Focus()
		case key.Matches(msg, m.keyMap.confirm):
			m.inInput = false
			m.input.Blur()
			m.updateKeybinds()
		case key.Matches(msg, m.keyMap.cancel):
			m.inInput = false
			m.input.Blur()
			m.updateKeybinds()
		case key.Matches(msg, m.keyMap.clear):
			m.input.SetValue("")
			m.updateKeybinds()
		case key.Matches(msg, m.keyMap.login):
			if !m.inNew {
				return m, m.loginCachedCmd(i)
			}
			return m, m.loginCmd
		case key.Matches(msg, m.keyMap.logout):
			return m, m.logoutCmd
		case key.Matches(msg, m.keyMap.open):
			return m, m.openAuthURLCmd
		case key.Matches(msg, m.keyMap.new):
			m.inNew = true
			m.updateKeybinds()
		case key.Matches(msg, m.keyMap.delete):
			return m, m.deleteUserCmd(i)
		}
	case NotificationMsg:
		return m, m.notifyCmd(string(msg))
	case NotificationTimeoutMsg:
		m.notification = ""
		return m, nil
	case error:
		return m, m.notifyCmd(msg.Error())
	}

	switch m.state {
	case Uninitialized:
		return m, func() tea.Msg {
			return errors.New("anilist model needs to be initialized")
		}
	case LoggedOut:
		if m.inNew {
			return m, m.updateInputCmd(msg)
		}
		return m, m.list.Update(msg)
	default: // and LoggedIn
		return m, nil
	}
}

// UpdateI is a convenience method to update the model in-place,
// to make it easier to handle.
func (m *Model) UpdateI(msg tea.Msg) tea.Cmd {
	self, cmd := m.Update(msg)
	m = self.(*Model)
	return cmd
}

// View implements tea.Model.
func (m *Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case Uninitialized:
		return "uninitialized"
	case LoggedIn:
		return m.viewLoggedIn()
	case LoggedOut:
		if m.inNew {
			return m.viewNewLogin()
		}
		return m.viewLoggedOut()
	default:
		return "unkown anilist login state"
	}
}
