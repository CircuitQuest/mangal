package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

func New(title, content string, color lipgloss.Color) base.State {
	viewport := viewport.New(0, 0)
	viewport.SetContent(content)
	viewport.Style = style.Normal.Base.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(color)

	_keyMap := newKeyMap(viewport.KeyMap)
	return &State{
		viewport: viewport,
		title:    title,
		keyMap:   _keyMap,
	}
}
