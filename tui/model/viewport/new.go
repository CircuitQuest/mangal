package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func New() *Model {
	v := viewport.New(0, 0)
	b := lipgloss.RoundedBorder()
	s := &Model{
		Model:                v,
		borderHorizontalSize: b.GetLeftSize() + b.GetRightSize(),
		borderVerticalSize:   b.GetTopSize() + b.GetBottomSize(),
		style:                lipgloss.NewStyle().BorderStyle(b),
		KeyMap:               newKeyMap(&v.KeyMap),
	}
	s.updateKeybinds()
	return s
}
