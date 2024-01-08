package manager

import (
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/provider/loader"
)

// TODO: need to provide options such as Mangadex specific options
func Loaders() ([]libmangal.ProviderLoader, error) {
	var loaders []libmangal.ProviderLoader

	mangoLoaders, err := loader.MangoLoaders()
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
