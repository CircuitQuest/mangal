package base

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Back goes back to the previous state if available.
func Back() tea.Msg {
	return BackMsg{}
}

// Notify sends a notification with the default time.Duration.
func Notify(message string) tea.Cmd {
	return func() tea.Msg {
		return NotificationMsg{Message: message}
	}
}

// NotifyWithDuration sends a notification with the given time.Duration.
func NotifyWithDuration(message string, duration time.Duration) tea.Cmd {
	return func() tea.Msg {
		return NotificationWithDurationMsg{
			NotificationMsg: NotificationMsg{
				Message: message,
			},
			Duration: duration,
		}
	}
}
