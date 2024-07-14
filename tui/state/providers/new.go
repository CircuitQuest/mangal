package providers

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/model/list"
	"github.com/zyedidia/generic/set"
)

func New(loaders []libmangal.ProviderLoader) *state {
	_keyMap := newKeyMap()
	extraInfo := false
	loaded := set.NewMapset[*item]()
	listWrapper := list.New(
		2,
		"provider", "providers",
		loaders,
		func(loader libmangal.ProviderLoader) _list.DefaultItem {
			return &item{
				loader:      loader,
				loadedItems: &loaded,
				extraInfo:   &extraInfo,
			}
		},
		&_keyMap)

	return &state{
		list:      listWrapper,
		loaded:    &loaded,
		extraInfo: &extraInfo,
		keyMap:    &_keyMap,
	}
}
