package chapsdownloading

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
)

type Options struct {
	DownloadChapter    func(ctx context.Context, chapter libmangal.Chapter) (libmangal.DownloadedChapter, error)
	OnDownloadFinished func(downChaps []libmangal.DownloadedChapter, succeed, failed []libmangal.Chapter) tea.Cmd
}
