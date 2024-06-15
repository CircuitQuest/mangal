package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/luevano/mangal/tui/base"
)

func New(title, content string, size base.Size) base.State {
	return &State{
		viewport: viewport.New(size.Width-2, size.Height-2), // -2 takes into account the border
		title:    title,
		content:  content,
		keyMap:   newKeyMap(),
	}
}
