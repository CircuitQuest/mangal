package viewport

import "github.com/luevano/mangal/tui/base"

func New(title, content string, size base.Size) *State {
	return &State{
		size:    size,
		title:   title,
		content: content,
		padding: base.Size{Width: 2, Height: 0},
		keyMap:  KeyMap{},
		styles:  DefaultStyles(),
	}
}
