package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
)

func New(title base.Title, content string, borderColor lipgloss.Color) *State {
	viewport := viewport.New(0, 0)
	viewport.SetContent(content)

	b := lipgloss.RoundedBorder()
	viewport.Style = style.Normal.Base.
		BorderStyle(b).
		BorderForeground(borderColor)

	s := &State{
		viewport:             viewport,
		title:                title,
		content:              content,
		keyMap:               newKeyMap(&viewport.KeyMap),
		borderHorizontalSize: b.GetLeftSize() + b.GetRightSize(),
		borderVerticalSize:   b.GetTopSize() + b.GetBottomSize(),
	}
	s.updateKeybinds()
	return s
}
