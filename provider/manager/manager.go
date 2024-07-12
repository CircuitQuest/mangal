package manager

import (
	"encoding/gob"

	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/provider/loader"
	mango "github.com/luevano/mangoprovider"
)

func Loaders() ([]libmangal.ProviderLoader, error) {
	// httpStoreProvider uses gob, types that will be stores must be registered
	gob.RegisterName("mango-manga", &mango.Manga{})
	gob.RegisterName("mango-volume", &mango.Volume{})
	gob.RegisterName("mango-chapter", &mango.Chapter{})
	gob.RegisterName("mango-page", &mango.Page{})
	gob.RegisterName("provider-metadata", &mangadata.Metadata{})
	gob.RegisterName("anilist-manga", &anilist.Manga{})

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
