package confirm

import (
	"github.com/luevano/mangal/tui/util"
)

func New(message string, onResponse OnResponseFunc) *State {
	return &State{
		message: message,
		keyMap: keyMap{
			yes: util.Bind("yes", "y", "enter"),
			no:  util.Bind("no", "n"),
		},
		onResponse: onResponse,
	}
}
