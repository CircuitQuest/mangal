package chapsdownloaded

import (
	"github.com/luevano/mangal/tui/util"
)

func New(options Options) *State {
	state := &State{
		options: options,
		keyMap: keyMap{
			open:  util.Bind("open directory", "o"),
			quit:  util.Bind("quit", "q"),
			retry: util.Bind("retry", "r"),
		},
	}
	state.keyMap.state = state
	return state
}
