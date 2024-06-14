package base

import (
	"context"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/zyedidia/generic/stack"
	"golang.org/x/term"
)

func New(state State,
	errState func(err error) State,
	logState func(title, content string, size Size) State,
) *model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 40
	}
	ctx, ctxCancel := context.WithCancel(context.Background())

	model := &model{
		state:     state,
		history:   stack.New[State](),
		ctx:       ctx,
		ctxCancel: ctxCancel,
		size: Size{
			Width:  width,
			Height: height,
		},
		styles:   defaultStyles(),
		keyMap:   newKeyMap(),
		help:     help.New(),
		errState: errState,
		logState: logState,
	}

	defer model.resize(model.stateSize())

	return model
}
