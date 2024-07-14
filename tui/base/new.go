package base

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/model/help"
	"github.com/luevano/mangal/tui/model/viewport"
)

func New(state State) *model {
	ctx, ctxCancel := context.WithCancel(context.Background())

	_styles := defaultStyles()
	_spinner := spinner.New(
		spinner.WithSpinner(_styles.loading.spinner),
		spinner.WithStyle(_styles.loading.spinnerStyle),
	)

	model := &model{
		state:                       state,
		viewport:                    viewport.New(),
		history:                     &history{},
		ctx:                         ctx,
		ctxCancel:                   ctxCancel,
		styles:                      _styles,
		keyMap:                      newKeyMap(),
		help:                        help.New(),
		spinner:                     _spinner,
		notificationDefaultDuration: time.Second,
		showBreadcrumbs:             config.TUI.ShowBreadcrumbs.Get(),
		showLoadingMessage:          true,
		showSubtitle:                true,
	}

	return model
}
