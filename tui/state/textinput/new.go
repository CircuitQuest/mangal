package textinput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/luevano/mangal/tui/util"
	"github.com/luevano/mangal/ui/icon"
)

func New(options Options) *State {
	if options.Title.Text == "" {
		options.Title.Text = "Search"
	}

	input := textinput.New()

	if options.Placeholder == "" {
		input.Placeholder = "Search..."
	} else {
		input.Placeholder = options.Placeholder
	}

	input.Prompt = fmt.Sprint(icon.Search, " ")

	return &State{
		options:   options,
		textinput: input,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
