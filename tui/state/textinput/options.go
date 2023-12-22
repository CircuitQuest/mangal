package textinput

import "github.com/luevano/mangal/tui/base"

type Options struct {
	Title        base.Title
	Prompt       string
	Placeholder  string
	Intermediate bool
	OnResponse   OnResponseFunc
}
