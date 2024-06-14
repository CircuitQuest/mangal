package loading

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/mangal/tui/base"
)

func New(message, subtitle string) *State {
	return &State{
		message:  message,
		subtitle: subtitle,
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
		keyMap:   base.NoKeyMap{},
	}
}
