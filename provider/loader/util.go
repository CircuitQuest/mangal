package loader

import (
	"path/filepath"
	"time"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/util/cache/bbolt"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
)

func cacheStore(dbName, bucketName string) (gokv.Store, error) {
	ttl, err := time.ParseDuration(config.Cache.TTL.Get())
	if err != nil {
		return nil, err
	}

	return bbolt.NewStore(bbolt.Options{
		TTL:        ttl,
		BucketName: bucketName,
		Path:       filepath.Join(path.CacheDir(), dbName+".db"),
		Codec:      encoding.Gob,
	})
}
