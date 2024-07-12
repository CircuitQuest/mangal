package metadata

import "github.com/luevano/libmangal/metadata"

func New(meta metadata.Metadata) *Model {
	m := &Model{
		styles: defaultStyles(),
	}
	m.SetMetadata(meta)
	return m
}
