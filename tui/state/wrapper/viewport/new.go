package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

func New(title base.Title, content string, borderColor lipgloss.Color) base.State {
	viewport := viewport.New(0, 0)
	viewport.SetContent(content)
	viewport.Style = style.Normal.Base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	return &State{
		viewport: viewport,
		title:    title,
		keyMap:   newKeyMap(viewport.KeyMap),
	}
}
