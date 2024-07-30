package cache

import (
	"encoding/gob"
	"errors"
	"path/filepath"
	"time"

	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/libmangal/metadata/myanimelist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/util/afs"
	"github.com/luevano/mangal/util/cache/bbolt"
	mango "github.com/luevano/mangoprovider"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"golang.org/x/oauth2"
)

const metadataDir = "metadata"

func init() {
	// CacheStore uses gob; types that will be stores must be registered first

	// manga data
	gob.RegisterName("mango-manga", &mango.Manga{})
	gob.RegisterName("mango-volume", &mango.Volume{})
	gob.RegisterName("mango-chapter", &mango.Chapter{})
	gob.RegisterName("mango-page", &mango.Page{})

	// metadatas
	gob.RegisterName("provider-metadata", &mangadata.Metadata{})
	gob.RegisterName("anilist-manga", &anilist.Manga{})
	gob.RegisterName("myanimelist-manga", &myanimelist.Manga{})

	// metadata users
	gob.RegisterName("anilist-user", &anilist.User{})
	gob.RegisterName("myanimelist-user", &myanimelist.User{})

	// oauth tokens
	gob.RegisterName("oauth-token", &oauth2.Token{})
}

// TODO: add subdir as parameter to avoid checking if the db is supposed to go into metadata dir
//
// CacheStore is the general cache store builder with a provided db name and bucket name.
//
// The TTL is decided based on the db name or bucket name.
func CacheStore(dbName, bucketName string) (gokv.Store, error) {
	// Dir
	dir := path.CacheDir()
	switch dbName {
	// metadata cache dbs into the subdir "metadata"
	case string(metadata.IDCodeAnilist),
		string(metadata.IDCodeMyAnimeList),
		string(metadata.IDCodeKitsu),
		string(metadata.IDCodeMangaUpdates),
		string(metadata.IDCodeAnimePlanet):
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
		return nil, errors.New(`can't use reserved bucket name "` + bbolt.TTLBucketName + `"`)
	case metadata.CacheBucketNameQueryToIDs:
		ttl = time.Hour * 24 * 2
	case metadata.CacheBucketNameTitleToID:
		// TODO: keeping the same old behavior, need to rework this;
		// the infinite TTL is implemented now
		ttl = time.Hour * 9999
	case metadata.CacheBucketNameIDToManga:
		ttl = time.Hour * 24 * 2
	case BucketNameAuthHistory,
		BucketNameAnilistAuthData,
		BucketNameMyAnimeListAuthData:
		// no expiry on the auth data, when the
		// tokens expire, delete the individual
		// auth data and prompt for re-authenticat
		ttl = 0
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
