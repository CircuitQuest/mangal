package providers

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
)

func loadProviderCmd(item *Item) tea.Cmd {
	return func() tea.Msg {
		return loadProviderMsg{
			item: item,
		}
	}
}

func searchMangasCmd(client *libmangal.Client) tea.Cmd {
	return func() tea.Msg {
		return searchMangasMsg{
			client: client,
		}
	}
}
