package base

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/model/help"
)

func New(state State,
	errState func(error) State,
	logState func() State,
) *model {
	ctx, ctxCancel := context.WithCancel(context.Background())

	_styles := defaultStyles()
	_spinner := spinner.New(
		spinner.WithSpinner(_styles.loading.spinner),
		spinner.WithStyle(_styles.loading.spinnerStyle),
	)

	model := &model{
		state:                       state,
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
		errState:                    errState,
		logState:                    logState,
	}

	return model
}
