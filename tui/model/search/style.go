package search

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type styles struct {
	renderSuggestionBox func(string) string

	suggestionsBox,
	suggestions,
	normalSuggestion,
	matchSuggestion lipgloss.Style
}

func defaultStyles() styles {
	// TODO: make width dynamic (this currently matches input width)
	sugs := lipgloss.NewStyle().Width(64).MaxWidth(64)
	sugsBox := lipgloss.NewStyle().
		MarginLeft(1). // alings with the search input box
		Border(lipgloss.RoundedBorder()).
		BorderTop(false).
		BorderForeground(color.Secondary)
	s := styles{
		suggestions:      sugs,
		normalSuggestion: style.Normal.Secondary,
		matchSuggestion:  style.Normal.Base,
		suggestionsBox:   sugsBox,
	}
	// to use the same values as the ones assigned, else
	// if these were to change, the function styling would be outdated
	s.renderSuggestionBox = func(str string) string {
		return s.suggestionsBox.Render(s.suggestions.Render(str))
	}

	return s
}
