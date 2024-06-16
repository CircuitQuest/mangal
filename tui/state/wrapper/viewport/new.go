package viewport

import "github.com/luevano/mangal/tui/base"

func New(title, content string) base.State {
	return &State{
		title:   title,
		content: content,
	}
}
