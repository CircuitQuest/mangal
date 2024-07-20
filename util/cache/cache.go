package cache

import (
	"encoding/gob"
	"errors"
	"path/filepath"
	"time"

	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/util/afs"
	"github.com/luevano/mangal/util/cache/bbolt"
	mango "github.com/luevano/mangoprovider"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
)

const metadataDir = "metadata"

func init() {
	// CacheStore uses gob; types that will be stores must be registered first
	gob.RegisterName("mango-manga", &mango.Manga{})
	gob.RegisterName("mango-volume", &mango.Volume{})
	gob.RegisterName("mango-chapter", &mango.Chapter{})
	gob.RegisterName("mango-page", &mango.Page{})
	gob.RegisterName("provider-metadata", &mangadata.Metadata{})
	gob.RegisterName("anilist-manga", &anilist.Manga{})
	gob.RegisterName("anilist-user", &anilist.User{})
}

func CacheStore(dbName, bucketName string) (gokv.Store, error) {
	// Dir
	dir := path.CacheDir()
	switch dbName {
	case anilist.CacheDBName:
		dir = filepath.Join(dir, metadataDir)
	}
	if err := afs.Afero.MkdirAll(dir, config.Download.ModeDir.Get()); err != nil {
		return nil, err
	}

	// TTL
	ttl, err := time.ParseDuration(config.Cache.TTL.Get())
	if err != nil {
		return nil, err
	}
	switch bucketName {
	case bbolt.TTLBucketName:
		return nil, errors.New(`can't use reserved bucket name "` + bucketName + `"`)
	case anilist.CacheBucketNameQueryToIDs:
		ttl = time.Hour * 24 * 2
	case anilist.CacheBucketNameTitleToID:
		// TODO: keeping the same old behavior, need to rework this;
		// the infinite TTL is implemented now
		ttl = time.Hour * 9999
	case anilist.CacheBucketNameIDToManga:
		ttl = time.Hour * 24 * 2
	case anilist.CacheBucketNameNameToAccessToken,
		anilist.CacheBucketNameNameToUser,
		BucketNameAnilistAuthHistory:
		ttl = time.Hour * 24 * 365 // access tokens last a year
	case BucketNameSearchHistory:
		ttl = 0 // no expiry
	}

	return bbolt.NewStore(bbolt.Options{
		TTL:        ttl,
		BucketName: bucketName,
		Path:       filepath.Join(dir, dbName+".db"),
		Codec:      encoding.Gob,
	})
}
