package base

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/mangal/theme/style"
	"github.com/zyedidia/generic/stack"
)

func New(state State,
	errState func(err error) State,
	logState func(title, content string) State,
) *model {
	ctx, ctxCancel := context.WithCancel(context.Background())

	_styles := defaultStyles()
	_help := help.New()
	_help.Styles = _styles.help
	_help.Ellipsis = Ellipsis
	_help.ShortSeparator = HelpKeySeparator
	_help.FullSeparator = HelpKeySeparator

	_spinner := spinner.New(
		spinner.WithSpinner(DotSpinner),
		spinner.WithStyle(style.Bold.Accent),
	)

	model := &model{
		state:                       state,
		history:                     stack.New[State](),
		ctx:                         ctx,
		ctxCancel:                   ctxCancel,
		styles:                      _styles,
		keyMap:                      newKeyMap(),
		help:                        _help,
		spinner:                     _spinner,
		notificationDefaultDuration: time.Second,
		showLoadingMessage:          true,
		showSubtitle:                true,
		errState:                    errState,
		logState:                    logState,
	}

	return model
}
