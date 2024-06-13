package providers

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New(loaders []libmangal.ProviderLoader) *State {
	listWrapper := list.New(util.NewList(
		3,
		"provider", "providers",
		loaders,
		func(loader libmangal.ProviderLoader) _list.DefaultItem {
			return Item{loader}
		},
	))

	return &State{
		providerLoaders: loaders,
		list:            listWrapper,
		keyMap: keyMap{
			info:    util.Bind("info", "i"),
			confirm: util.Bind("confirm", "enter"),
			list:    listWrapper.KeyMap(),
		},
	}
}
