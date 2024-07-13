package metadata

import (
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
	"github.com/luevano/mangal/tui/util"
)

func New(meta metadata.Metadata) *State {
	s := util.MetaIDStyle(meta.ID())
	p := meta.ID().Code
	if p != "" {
		p = "[" + p + "] "
	}
	title := base.Title{
		Text:       p + s.Prefix + " Metadata",
		Background: s.Color,
		Foreground: color.Bright,
	}
	return &State{
		viewport: viewport.New(title, "", s.Color),
		meta:     meta,
		styles:   defaultStyles(s),
	}
}
