package chapsdownloading

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
)

type Options struct {
	DownloadChapter    func(ctx context.Context, chapter mangadata.Chapter) (*metadata.DownloadedChapter, error)
	OnDownloadFinished func(downChaps []*metadata.DownloadedChapter, succeed, failed []mangadata.Chapter) tea.Cmd
}
