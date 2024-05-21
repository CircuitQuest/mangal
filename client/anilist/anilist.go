package anilist

import (
	"log"
	"path/filepath"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/util/afs"
	"github.com/luevano/mangal/util/cache/bbolt"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
)

var Anilist = newAnilist()

func newAnilist() *libmangal.Anilist {
	newPersistentStore := func(name string, ttl time.Duration) (gokv.Store, error) {
		dir := filepath.Join(path.CacheDir(), "anilist")
		if err := afs.Afero.MkdirAll(dir, config.Config.Download.ModeDir.Get()); err != nil {
			return nil, err
		}

		return bbolt.NewStore(bbolt.Options{
			TTL:        ttl,
			BucketName: name,
			Path:       filepath.Join(dir, name+".db"),
			Codec:      encoding.Gob,
		})
	}

	anilistOptions := libmangal.DefaultAnilistOptions()

	var err error
	anilistOptions.QueryToIDsStore, err = newPersistentStore("query-to-id", time.Hour*24*2)
	if err != nil {
		log.Fatal(err)
	}

	anilistOptions.IDToMangaStore, err = newPersistentStore("id-to-manga", time.Hour*24*2)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: make infinite ttl
	anilistOptions.TitleToIDStore, err = newPersistentStore("title-to-id", time.Hour*9999)
	if err != nil {
		log.Fatal(err)
	}

	anilistOptions.AccessTokenStore, err = newPersistentStore("access-token", time.Hour*24*30)
	if err != nil {
		log.Fatal(err)
	}

	anilist := libmangal.NewAnilist(anilistOptions)
	return &anilist
}
