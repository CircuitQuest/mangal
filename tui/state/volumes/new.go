package volumes

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(client *libmangal.Client, manga mangadata.Manga, volumes []mangadata.Volume) *state {
	keyMap := newKeyMap()
	listWrapper := list.New(
		1,
		"volume", "volumes",
		volumes,
		func(volume mangadata.Volume) _list.DefaultItem {
			return &item{volume}
		},
		keyMap)

	return &state{
		list:    listWrapper,
		volumes: volumes,
		manga:   manga,
		client:  client,
		keyMap:  keyMap,
	}
}
