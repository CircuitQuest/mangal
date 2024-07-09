package providers

import tea "github.com/charmbracelet/bubbletea"

func loadProviderCmd(item *item) tea.Cmd {
	return func() tea.Msg {
		return loadProviderMsg{
			item: item,
		}
	}
}
