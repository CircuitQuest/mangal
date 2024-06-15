package errorstate

import "github.com/luevano/mangal/tui/base"

func New(err error, size base.Size) base.State {
	return &State{
		error:  err,
		size:   size,
		keyMap: newKeyMap(),
	}
}
