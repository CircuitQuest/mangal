package viewport

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

// styles for the viewport.
type styles struct {
	Title,
	Info,
	Line,
	Viewport lipgloss.Style
	Content        func(int) lipgloss.Style
	ContentWrapper func(int, int) lipgloss.Style
}

// defaultStyles get a sensible default viewport style.
func defaultStyles() styles {
	borderRight := lipgloss.RoundedBorder()
	borderRight.Right = "├"
	borderRight.BottomLeft = "├"
	borderLeft := lipgloss.RoundedBorder()
	borderLeft.Left = "┤"
	borderLeft.TopRight = "┤"

	boxStyle := lipgloss.NewStyle().BorderForeground(color.Viewport).Padding(0, 1)
	return styles{
		Title:    boxStyle.BorderStyle(borderRight),
		Info:     boxStyle.BorderStyle(borderLeft),
		Line:     style.Normal.Viewport,
		Viewport: boxStyle.BorderStyle(lipgloss.RoundedBorder()).BorderLeft(true).BorderRight(true),
		// TODO: use style.Trim?
		Content: func(maxSize int) lipgloss.Style {
			return style.Normal.Base.Width(maxSize)
		},
		ContentWrapper: func(height, width int) lipgloss.Style {
			return style.Normal.Base.Padding(height, width)
		},
	}
}
