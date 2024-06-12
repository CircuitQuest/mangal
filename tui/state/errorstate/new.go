package errorstate

import "github.com/luevano/mangal/tui/util"

func New(err error) *State {
	return &State{
		error: err,
		keyMap: keyMap{
			quit:      util.Bind("quit", "q"),
			copyError: util.Bind("copy error", "c"),
		},
	}
}
