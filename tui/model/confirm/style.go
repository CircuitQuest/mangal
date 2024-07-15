package confirm

import "github.com/charmbracelet/lipgloss"

type styles struct {
	title,
	message lipgloss.Style
}

func (s styles) updateWidth(width int) {
	s.message = s.message.Width(width).MaxWidth(width)
}

func defaultStyles(width int, color lipgloss.Color) styles {
	return styles{
		title:   lipgloss.NewStyle().Padding(0, 1).Background(color),
		message: lipgloss.NewStyle().Width(width).MaxWidth(width),
	}
}
