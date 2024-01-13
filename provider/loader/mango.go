package loader

import (
	"net/http"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangoprovider"
	"github.com/luevano/mangoprovider/mango"
)

func MangoLoaders(options Options) ([]libmangal.ProviderLoader, error) {
	o := mango.Options{
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		HTTPStoreProvider: httpStoreProvider,
		Filter: mango.Filter{
			NSFW:                    options.NSFW,
			Language:                options.Language,
			MangaDexDataSaver:       options.MangaDexDataSaver,
			TitleChapterNumber:      options.TitleChapterNumber,
			AvoidDuplicateChapters:  options.AvoidDuplicateChapters,
			ShowUnavailableChapters: options.ShowUnavailableChapters,
		},
	}

	return mangoprovider.Loaders(o)
}
