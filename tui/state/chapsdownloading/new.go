package chapsdownloading

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/theme/style"
)

func New(chapters []mangadata.Chapter, options Options) *State {
	return &State{
		options:  options,
		chapters: chapters,
		message:  "Preparing...",
		progress: progress.New(),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Dot),
			spinner.WithStyle(style.Bold.Accent),
		),
		keyMap: KeyMap{},
	}
}
