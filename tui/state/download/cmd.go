package download

import tea "github.com/charmbracelet/bubbletea"

func startDownloadCmd() tea.Msg {
	return startDownloadMsg{}
}

func downloadFailedCmd() tea.Msg {
	return downloadFailedMsg{}
}

func nextChapterCmd() tea.Msg {
	return nextChapterMsg{}
}
