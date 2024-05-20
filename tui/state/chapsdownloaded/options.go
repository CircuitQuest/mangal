package chapsdownloaded

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
)

type Options struct {
	Succeed, Failed  []libmangal.Chapter
	SucceedDownloads []libmangal.DownloadedChapter
	DownloadChapters func(chapters []libmangal.Chapter) tea.Cmd
}
