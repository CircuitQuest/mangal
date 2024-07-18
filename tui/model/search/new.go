package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/tui/base"
)

func New(placeholder, initialQuery string, width, maxSuggestions int) *Model {
	input := textinput.New()
	input.Placeholder = placeholder
	if input.Placeholder == "" {
		input.Placeholder = "Search..."
	}

	input.ShowSuggestions = true
	input.Width = width
	input.Prompt = icon.Search.Colored() + " "

	initState := Unsearched
	query := strings.TrimSpace(initialQuery)
	if query != "" {
		input.SetValue(query)
		initState = Searched
	}

	return &Model{
		input:    input,
		state:    initState,
		query:    query,
		maxWidth: width,
		size:     base.Size{Width: width, Height: maxSuggestions},
		styles:   defaultStyles(),
		keyMap:   newKeyMap(),
	}
}
