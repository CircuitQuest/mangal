package anilist

import (
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/util/afs"
	"github.com/luevano/mangal/util/cache/bbolt"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
)

var Anilist = newAnilist()

func newAnilist() *anilist.Anilist {
	options := anilist.DefaultOptions()
	options.CacheStore = cacheStore

	ani, err := anilist.NewAnilist(options)
	if err != nil {
		log.Fatal(err)
	}
	return ani
}

func cacheStore(dbName, bucketName string) (gokv.Store, error) {
	dir := filepath.Join(path.CacheDir(), "metadata")
	if err := afs.Afero.MkdirAll(dir, config.Download.ModeDir.Get()); err != nil {
		return nil, err
	}

	ttl := time.Hour * 24 * 2
	switch bucketName {
	case bbolt.TTLBucketName:
		return nil, errors.New(`can't use reserved bucket name "` + bucketName + `"`)
	case anilist.CacheBucketNameTitleToID:
		// TODO: keeping the same old behavior, need to rework this;
		// the infinite TTL is implemented now
		ttl = time.Hour * 9999
	case anilist.CacheBucketNameAccessToken:
		ttl = time.Hour * 24 * 30
	}

	return bbolt.NewStore(bbolt.Options{
		TTL:        ttl,
		BucketName: bucketName,
		Path:       filepath.Join(dir, dbName+".db"),
		Codec:      encoding.Gob,
	})
}
