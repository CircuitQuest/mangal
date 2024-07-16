package viewport

import tea "github.com/charmbracelet/bubbletea"

func notificationCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return NotificationMsg{
			Message: message,
		}
	}
}

func backCmd() tea.Msg {
	return BackMsg{}
}
