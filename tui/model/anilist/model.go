package anilist

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/metadata/anilist"
)

var authURL = "https://anilist.co/api/v2/oauth/authorize?client_id=" +
	"%s" +
	"&response_type=code&redirect_uri=https://anilist.co/api/v2/oauth/pin"

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
	help    help.Model
	anilist *anilist.Anilist

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
	inInput bool

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
	if m.anilist.IsAuthorized() {
		m.state = LoggedIn
	} else {
		m.state = LoggedOut
	}
	m.updateKeybinds()
	return textinput.Blink
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.quit):
			m.quitting = true
			return m, tea.Quit
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
		case key.Matches(msg, m.keyMap.login):
			return m, m.loginCmd
		case key.Matches(msg, m.keyMap.open):
			return m, m.openAuthURLCmd
		case key.Matches(msg, m.keyMap.logout):
			return m, m.logoutCmd
		}
	case NotificationMsg:
		return m, m.notifyCmd(string(msg))
	case NotificationTimeoutMsg:
		m.notification = ""
		return m, nil
	}
	if m.Uninitialized() || !m.inInput {
		return m, nil
	}

	return m, m.updateInputCmd(msg)
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
	if m.Uninitialized() {
		return "Uninitialized"
	}

	view := "Already logged in"
	if m.LoggedOut() {
		view = m.viewLogin()
	}

	if !m.standalone {
		return view
	}

	msg := m.styles.notification.Render(m.notification)
	view = lipgloss.JoinVertical(lipgloss.Left, m.title, msg, view, m.help.View(&m.keyMap))
	return m.styles.view.Render(view)
}

func (m *Model) DisableQuitKeybindings() {
	m.keyMap.quit.SetEnabled(false)
}
