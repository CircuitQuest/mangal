package metadata

import (
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/color"
	_style "github.com/luevano/mangal/theme/style"
)

func New(meta metadata.Metadata) *Model {
	m := &Model{
		style: _style.Normal.Base.Padding(0, 1).Foreground(color.Bright),
	}
	m.SetMetadata(meta)
	return m
}
