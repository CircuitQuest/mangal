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

	model := &model{
		state:                       state,
		history:                     stack.New[State](),
		ctx:                         ctx,
		ctxCancel:                   ctxCancel,
		styles:                      defaultStyles(),
		keyMap:                      newKeyMap(),
		help:                        help.New(),
		notificationDefaultDuration: time.Second,
		errState:                    errState,
		logState:                    logState,
	}

	return model
}
