package path

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) notifyCmd(message string) tea.Cmd {
	m.notification = message
	return tea.Tick(m.notificationDuration, func(t time.Time) tea.Msg {
		return NotificationTimeoutMsg{}
	})
}
