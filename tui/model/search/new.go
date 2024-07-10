package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
)

func New(placeholder, initialQuery string) *Model {
	input := textinput.New()
	input.Placeholder = placeholder
	if input.Placeholder == "" {
		input.Placeholder = "Search..."
	}

	input.Prompt = icon.Search.String() + " "
	input.PromptStyle = style.Bold.Warning
	input.CharLimit = 64

	initState := Unsearched
	query := strings.TrimSpace(initialQuery)
	if query != "" {
		input.SetValue(query)
		initState = Searched
	}

	return &Model{
		input:  input,
		state:  initState,
		query:  query,
		style:  lipgloss.NewStyle().Padding(0, 1, 1, 1),
		keyMap: newKeyMap(),
	}
}
