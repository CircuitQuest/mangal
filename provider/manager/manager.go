package manager

import (
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/provider/loader"
)

func Loaders(options loader.Options) ([]libmangal.ProviderLoader, error) {
	var loaders []libmangal.ProviderLoader

	mangoLoaders, err := loader.MangoLoaders(options)
	if err != nil {
		return nil, err
	}
	loaders = append(loaders, mangoLoaders...)

	luaLoaders, err := loader.LuaLoaders()
	if err != nil {
		return nil, err
	}
	loaders = append(loaders, luaLoaders...)

	return loaders, nil
}
