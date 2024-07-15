package confirm

import tea "github.com/charmbracelet/bubbletea"

func yesCmd() tea.Msg {
	return YesMsg{}
}

func noCmd() tea.Msg {
	return NoMsg{}
}
