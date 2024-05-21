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

func MangoLoaders(options Options) ([]libmangal.ProviderLoader, error) {
	// httpStoreProvider uses gob, it needs to register custom types
	gob.Register(&mango.Manga{})
	gob.Register(&mango.Volume{})
	gob.Register(&mango.Chapter{})
	gob.Register(&mango.Page{})

	o := mango.Options{
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		UserAgent:   config.Config.Download.UserAgent.Get(),
		HTTPStore:   httpStore,
		Parallelism: options.Parallelism,
		Headless: mango.Headless{
			UseFlaresolverr: options.HeadlessUseFlaresolverr,
			FlaresolverrURL: options.HeadlessFlaresolverrURL,
		},
		Filter: mango.Filter{
			NSFW:                    options.NSFW,
			Language:                options.Language,
			MangaPlusQuality:        options.MangaPlusQuality,
			MangaDexDataSaver:       options.MangaDexDataSaver,
			TitleChapterNumber:      options.TitleChapterNumber,
			AvoidDuplicateChapters:  options.AvoidDuplicateChapters,
			ShowUnavailableChapters: options.ShowUnavailableChapters,
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
