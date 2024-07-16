package download

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/viewport"
)

func New(client *libmangal.Client, chaptersToDownload []mangadata.Chapter, options libmangal.DownloadOptions) *state {
	c := make(chapters, len(chaptersToDownload))
	for i, ch := range chaptersToDownload {
		c[i] = &chapter{
			chapter: ch,
			state:   cSToDownload,
		}
	}

	_styles := defaultStyles()
	sep := _styles.sep.Render(icon.Separator.Raw())

	_viewport := viewport.New()
	_keyMap := newKeyMap(&_viewport.KeyMap)
	return &state{
		progress: progress.New(),
		spinner: spinner.New(
			spinner.WithSpinner(base.DotSpinner),
			spinner.WithStyle(style.Normal.Accent),
		),
		viewport:    _viewport,
		client:      client,
		chapters:    c,
		options:     options,
		downloading: dSUninitialized,
		sep:         sep,
		message:     "Preparing...",
		styles:      _styles,
		keyMap:      _keyMap,
	}
}
