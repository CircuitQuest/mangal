package loader

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	mango "github.com/luevano/mangoprovider"
	"github.com/luevano/mangoprovider/apis"
	"github.com/luevano/mangoprovider/scrapers"
)

func MangoLoaders() ([]libmangal.ProviderLoader, error) {
	// httpStoreProvider uses gob, it needs to register custom types
	gob.Register(&mango.Manga{})
	gob.Register(&mango.Volume{})
	gob.Register(&mango.Chapter{})
	gob.Register(&mango.Page{})

	// Generates overall default options then overriding as necessary
	o := mango.DefaultOptions()

	o.HTTPClient.Timeout = time.Minute
	o.UserAgent = config.Download.UserAgent.Get()
	o.HTTPStore = httpStore
	o.Parallelism = config.Providers.Parallelism.Get()
	o.Filter.NSFW = config.Providers.Filter.NSFW.Get()
	o.Filter.Language = config.Providers.Filter.Language.Get()
	o.Filter.TitleChapterNumber = config.Providers.Filter.TitleChapterNumber.Get()
	o.Filter.AvoidDuplicateChapters = config.Providers.Filter.AvoidDuplicateChapters.Get()
	o.Filter.ShowUnavailableChapters = config.Providers.Filter.ShowUnavailableChapters.Get()
	o.Headless.UseFlaresolverr = config.Providers.Headless.UseFlaresolverr.Get()
	o.Headless.FlaresolverrURL = config.Providers.Headless.FlaresolverrURL.Get()
	o.MangaDex.DataSaver = config.Providers.MangaDex.DataSaver.Get()
	o.MangaPlus.Quality = config.Providers.MangaPlus.Quality.Get()
	o.MangaPlus.OSVersion = config.Providers.MangaPlus.OSVersion.Get()
	o.MangaPlus.AppVersion = config.Providers.MangaPlus.AppVersion.Get()
	// AndroidID is the only config defaulted to empty string,
	// mangoprovider generates a random one when using default options
	if androidID := config.Providers.MangaPlus.AndroidID.Get(); androidID != "" {
		o.MangaPlus.AndroidID = androidID
	}

	var loaders []libmangal.ProviderLoader
	loaders = append(loaders, apis.Loaders(o)...)
	loaders = append(loaders, scrapers.Loaders(o)...)

	for _, loader := range loaders {
		if loader == nil {
			// TODO: need to provide more info
			return nil, fmt.Errorf("failed while loading providers")
		}
	}

	return loaders, nil
}
