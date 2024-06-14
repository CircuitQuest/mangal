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
) *Model {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 40
	}

	model := &Model{
		state:   state,
		history: stack.New[State](),
		size: Size{
			Width:  width,
			Height: height,
		},
		keyMap:   newKeyMap(),
		help:     help.New(),
		styles:   DefaultStyles(),
		errState: errState,
		logState: logState,
	}

	defer model.resize(model.StateSize())

	model.ctx, model.ctxCancel = context.WithCancel(context.Background())

	return model
}
