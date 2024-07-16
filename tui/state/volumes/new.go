package volumes

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/luevano/mangal/tui/model/metadata"
)

func New(client *libmangal.Client, manga mangadata.Manga, volumes []mangadata.Volume) *state {
	listWrapper := list.New(
		1,
		"volume", "volumes",
		volumes,
		func(volume mangadata.Volume) _list.DefaultItem {
			return &item{volume}
		},
	)

	return &state{
		list:    listWrapper,
		meta:    metadata.New(manga.Metadata()),
		volumes: volumes,
		manga:   manga,
		client:  client,
		keyMap:  newKeyMap(),
	}
}
