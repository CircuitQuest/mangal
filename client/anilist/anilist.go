package anilist

import (
	"log"

	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/util/cache"
)

var Anilist = newAnilist()

func newAnilist() *anilist.Anilist {
	options := anilist.DefaultOptions()
	options.CacheStore = cache.CacheStore

	ani, err := anilist.NewAnilist(options)
	if err != nil {
		log.Fatal(err)
	}
	return ani
}
