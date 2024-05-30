package chapsdownloaded

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
)

type Options struct {
	Succeed, Failed  []mangadata.Chapter
	SucceedDownloads []*metadata.DownloadedChapter
	DownloadChapters func(chapters []mangadata.Chapter) tea.Cmd
}
