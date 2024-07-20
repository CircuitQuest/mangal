package anilist

import (
	"log"

	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/util/cache"
)

var anilist_ = newAnilist()

func Anilist() *anilist.Anilist {
	return anilist_
}

func newAnilist() *anilist.Anilist {
	options := anilist.DefaultOptions()
	options.CacheStore = cache.CacheStore

	// Authenticate with the last authenticated user if existent
	var userHistory cache.UserHistory
	_, _ = cache.GetAnilistAuthHistory(&userHistory)
	options.Username = userHistory.Last() // could be empty, which is fine

	ani, err := anilist.NewAnilist(options)
	if err != nil {
		log.Fatal(err)
	}
	return ani
}
