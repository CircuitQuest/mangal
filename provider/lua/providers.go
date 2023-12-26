package lua

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/cache/bbolt"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/libmangal"
	"github.com/luevano/luaprovider"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
)

const (
	mainLua = "main.lua"
)

func NewLoader(info libmangal.ProviderInfo, dir string) (libmangal.ProviderLoader, error) {
	providerMainFilePath := filepath.Join(dir, mainLua)
	exists, err := afs.Afero.Exists(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("%s is missing", providerMainFilePath)
	}

	providerMainFileContents, err := afs.Afero.ReadFile(providerMainFilePath)
	if err != nil {
		return nil, err
	}

	ttl, err := time.ParseDuration(config.Config.Providers.Cache.TTL.Get())
	if err != nil {
		log.Fatal(err)
	}

	options := luaprovider.Options{
		PackagePaths: []string{dir},
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		HTTPStoreProvider: func() (gokv.Store, error) {
			return bbolt.NewStore(bbolt.Options{
				TTL:        ttl,
				BucketName: info.Name,
				Path:       filepath.Join(path.CacheDir(), info.Name+".db"),
				Codec:      encoding.Gob,
			})
		},
	}

	return luaprovider.NewLoader(providerMainFileContents, info, options)
}
