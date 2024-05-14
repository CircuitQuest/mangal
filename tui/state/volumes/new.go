package volumes

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
	"github.com/luevano/libmangal"
)

func New(client *libmangal.Client, manga *libmangal.Manga, volumes []*libmangal.Volume) *State {
	listWrapper := listwrapper.New(util.NewList(
		1,
		"volume", "volumes",
		volumes,
		func(volume *libmangal.Volume) list.DefaultItem {
			return Item{volume}
		},
	))

	return &State{
		manga:   manga,
		client:  client,
		volumes: volumes,
		list:    listWrapper,
		keyMap: KeyMap{
			Confirm: util.Bind("confirm", "enter"),
		},
	}
}
