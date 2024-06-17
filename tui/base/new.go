package base

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/zyedidia/generic/stack"
)

func New(state State,
	errState func(err error) State,
	logState func(title, content string) State,
) *model {
	ctx, ctxCancel := context.WithCancel(context.Background())

	_styles := defaultStyles()
	_help := help.New()

	_help.Styles.ShortKey = _styles.helpKey
	_help.Styles.ShortSeparator = _styles.helpSep
	_help.Styles.FullKey = _styles.helpKey
	_help.Styles.FullSeparator = _styles.helpSep

	model := &model{
		state:                       state,
		history:                     stack.New[State](),
		ctx:                         ctx,
		ctxCancel:                   ctxCancel,
		styles:                      _styles,
		keyMap:                      newKeyMap(),
		help:                        _help,
		notificationDefaultDuration: time.Second,
		errState:                    errState,
		logState:                    logState,
	}

	return model
}
