package errorstate

import "github.com/luevano/mangal/tui/base"

func New(err error) base.State {
	return &State{
		error:  err,
		keyMap: newKeyMap(),
	}
}
