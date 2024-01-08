package loader

import (
	"net/http"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangoprovider"
	"github.com/luevano/mangoprovider/mango"
)

func MangoLoaders() ([]libmangal.ProviderLoader, error) {
	options := mango.Options{
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		HTTPStoreProvider: httpStoreProvider,
	}

	return mangoprovider.Loaders(options)
}
