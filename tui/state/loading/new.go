package loading

import "github.com/charmbracelet/bubbles/spinner"

func New(title, message string) *State {
	return &State{
		title:   title,
		message: message,
		spinner: spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}
