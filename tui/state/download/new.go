package download

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

func New(client *libmangal.Client, chapters []mangadata.Chapter, options libmangal.DownloadOptions) *state {
	return &state{
		client:   client,
		chapters: chapters,
		options:  options,
		message:  "Preparing...",
		progress: progress.New(),
		spinner: spinner.New(
			spinner.WithSpinner(base.DotSpinner),
			spinner.WithStyle(style.Normal.Accent),
		),
		keyMap: newKeyMap(),
	}
}
