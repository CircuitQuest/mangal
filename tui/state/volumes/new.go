package volumes

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New(client *libmangal.Client, manga mangadata.Manga, volumes []mangadata.Volume) *State {
	listWrapper := list.New(util.NewList(
		1,
		"volume", "volumes",
		volumes,
		func(volume mangadata.Volume) _list.DefaultItem {
			return Item{volume}
		},
	))

	return &State{
		list:    listWrapper,
		volumes: volumes,
		manga:   manga,
		client:  client,
		keyMap: keyMap{
			confirm: util.Bind("confirm", "enter"),
		},
	}
}
