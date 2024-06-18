package providers

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/zyedidia/generic/set"
)

func New(loaders []libmangal.ProviderLoader) *State {
	keyMap := newKeyMap()
	extraInfo := false
	loaded := set.NewMapset[*Item]()
	listWrapper := list.New(
		2,
		"provider", "providers",
		loaders,
		func(loader libmangal.ProviderLoader) _list.DefaultItem {
			return &Item{
				ProviderLoader: loader,
				loadedItems:    &loaded,
				extraInfo:      &extraInfo,
			}
		},
		keyMap)

	return &State{
		list:      listWrapper,
		loaded:    &loaded,
		extraInfo: &extraInfo,
		keyMap:    keyMap,
	}
}
