package providers

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
	"github.com/mangalorg/libmangal"
)

func New(loaders []libmangal.ProviderLoader) *State {
	listWrapper := listwrapper.New(util.NewList(
		4,
		"provider", "providers",
		loaders,
		func(loader libmangal.ProviderLoader) list.DefaultItem {
			return Item{loader}
		},
	))

	return &State{
		providersLoaders: loaders,
		list:             listWrapper,
		keyMap: KeyMap{
			info:    util.Bind("info", "i"),
			confirm: util.Bind("confirm", "enter"),
			list:    listWrapper.GetKeyMap(),
		},
	}
}
