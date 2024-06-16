package providers

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
)

func New(loaders []libmangal.ProviderLoader) *State {
	keyMap := newKeyMap()
	listWrapper := list.New(util.NewList(
		3,
		"provider", "providers",
		loaders,
		func(loader libmangal.ProviderLoader) _list.DefaultItem {
			return Item{loader}
		},
	), keyMap)

	return &State{
		providerLoaders: loaders,
		list:            listWrapper,
		keyMap:          keyMap,
	}
}
