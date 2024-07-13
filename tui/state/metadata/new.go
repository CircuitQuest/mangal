package metadata

import (
	"github.com/luevano/mangal/tui/model/metadata"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

func New(meta *metadata.Model) *State {
	return &State{
		viewport: viewport.New("Metadata", "", meta.Style().Color),
		meta:     meta,
	}
}
