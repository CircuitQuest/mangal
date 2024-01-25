package viewport

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/ui/color"
)

// Styles for the viewport.
type Styles struct {
	Title,
	Info,
	Line,
	Viewport lipgloss.Style
	Content        func(int) lipgloss.Style
	ContentWrapper func(int, int) lipgloss.Style
}

// DefaultStyles get a sensible default viewport style.
func DefaultStyles() Styles {
	borderRight := lipgloss.RoundedBorder()
	borderRight.Right = "├"
	borderRight.BottomLeft = "├"
	borderLeft := lipgloss.RoundedBorder()
	borderLeft.Left = "┤"
	borderLeft.TopRight = "┤"

	boxStyle := lipgloss.NewStyle().BorderForeground(color.Viewport).Padding(0, 1)
	return Styles{
		Title: boxStyle.Copy().
			BorderStyle(borderRight),
		Info: boxStyle.Copy().
			BorderStyle(borderLeft),
		Line: lipgloss.
			NewStyle().
			Foreground(color.Viewport),
		Viewport: boxStyle.Copy().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderLeft(true).
			BorderRight(true),
		Content: func(maxSize int) lipgloss.Style {
			return lipgloss.
				NewStyle().
				Width(maxSize)
		},
		ContentWrapper: func(height, width int) lipgloss.Style {
			return lipgloss.
				NewStyle().
				Padding(height, width)
		},
	}
}
