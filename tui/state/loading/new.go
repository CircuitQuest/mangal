package loading

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/luevano/mangal/tui/model"
)

func New(message, subtitle string) *State {
	return &State{
		message:  message,
		subtitle: subtitle,
		spinner:  spinner.New(spinner.WithSpinner(spinner.Dot)),
		keyMap:   model.NoKeyMap{},
	}
}
