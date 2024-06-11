package loader

import (
	"encoding/gob"
	"fmt"
	"net/http"
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

	o := mango.Options{
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		UserAgent:   config.Download.UserAgent.Get(),
		HTTPStore:   httpStore,
		Parallelism: config.Providers.Parallelism.Get(),
		Headless: mango.Headless{
			UseFlaresolverr: config.Providers.Headless.UseFlaresolverr.Get(),
			FlaresolverrURL: config.Providers.Headless.FlaresolverrURL.Get(),
		},
		Filter: mango.Filter{
			NSFW:                    config.Providers.Filter.NSFW.Get(),
			Language:                config.Providers.Filter.Language.Get(),
			MangaPlusQuality:        config.Providers.Filter.MangaPlusQuality.Get(),
			MangaDexDataSaver:       config.Providers.Filter.MangaDexDataSaver.Get(),
			TitleChapterNumber:      config.Providers.Filter.TitleChapterNumber.Get(),
			AvoidDuplicateChapters:  config.Providers.Filter.AvoidDuplicateChapters.Get(),
			ShowUnavailableChapters: config.Providers.Filter.ShowUnavailableChapters.Get(),
		},
		// These will only be set upstream if they're non-empty
		MangaPlus: mango.MangaPlusOptions{
			OSVersion:  config.Providers.MangaPlus.OSVersion.Get(),
			AppVersion: config.Providers.MangaPlus.AppVersion.Get(),
			AndroidID:  config.Providers.MangaPlus.AndroidID.Get(),
		},
	}
	var loaders []libmangal.ProviderLoader
	loaders = append(loaders, apis.Loaders(o)...)
	loaders = append(loaders, scrapers.Loaders(o)...)

	for _, loader := range loaders {
		if loader == nil {
			return nil, fmt.Errorf("failed while loading providers")
		}
	}

	return loaders, nil
}
