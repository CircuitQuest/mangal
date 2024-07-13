package metadata

import (
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
	"github.com/luevano/mangal/tui/util"
)

func New(meta metadata.Metadata) *State {
	metaStyle := util.MetaIDStyle(meta.ID())
	p := meta.ID().Code
	if p != "" {
		p = "[" + p + "] "
	}
	title := base.Title{
		Text:       p + metaStyle.Prefix + " Metadata",
		Background: metaStyle.Color,
		Foreground: color.Bright,
	}

	s := &State{
		meta:   meta,
		styles: defaultStyles(metaStyle),
	}
	viewport := viewport.New(title, s.renderMetadata(), metaStyle.Color)
	s.viewport = viewport
	return s
}
