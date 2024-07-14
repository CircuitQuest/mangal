package metadata

import "github.com/luevano/libmangal/metadata"

func New(meta metadata.Metadata) *Model {
	metaStyle := metaIDStyle(meta.ID())
	return &Model{
		meta:      meta,
		metaStyle: metaStyle,
		styles:    defaultStyles(metaStyle),
	}
}
