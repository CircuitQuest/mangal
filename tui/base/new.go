package base

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/config"
)

func New(state State,
	errState func(err error) State,
	logState func(title, content string, color lipgloss.Color) State,
) *model {
	ctx, ctxCancel := context.WithCancel(context.Background())

	_styles := defaultStyles()
	_help := help.New()
	_help.Styles = _styles.help
	_help.Ellipsis = Ellipsis
	_help.ShortSeparator = HelpKeySeparator
	_help.FullSeparator = HelpKeySeparator

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
		help:                        _help,
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
