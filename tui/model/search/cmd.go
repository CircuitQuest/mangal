package search

import tea "github.com/charmbracelet/bubbletea"

func SearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		return SearchMsg(query)
	}
}

func SearchCancelCmd() tea.Msg {
	return SearchCancelMsg{}
}
