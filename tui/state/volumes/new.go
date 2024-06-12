package volumes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
)

func New(client *libmangal.Client, manga *mangadata.Manga, volumes []*mangadata.Volume) *State {
	listWrapper := listwrapper.New(util.NewList(
		1,
		"volume", "volumes",
		volumes,
		func(volume *mangadata.Volume) list.DefaultItem {
			return Item{volume}
		},
	))

	return &State{
		manga:   manga,
		client:  client,
		volumes: volumes,
		list:    listWrapper,
		keyMap: keyMap{
			confirm: util.Bind("confirm", "enter"),
		},
	}
}
