package base

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Back goes back to the previous state if available.
func Back() tea.Msg {
	return BackMsg{
		Steps: 1,
	}
}

// BackN goes back N states back or up to the initial State (home).
func BackN(steps int) tea.Cmd {
	if steps < 0 {
		panic("steps < 0 in BackN")
	}

	return func() tea.Msg {
		return BackMsg{Steps: steps}
	}
}

// BackToHome goes back to the initial State (home).
func BackToHome() tea.Msg {
	return BackToHomeMsg{}
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

// Loading will display a loading message with a spinner.
func Loading(message string) tea.Cmd {
	return func() tea.Msg {
		return LoadingMsg{
			Message: message,
		}
	}
}

// Loaded will stop the loading message.
func Loaded() tea.Msg {
	return LoadingMsg{
		Message: "",
	}
}
