package anilist

import (
	"github.com/metafates/mangal/cache"
	"github.com/metafates/mangal/constant"
	"github.com/samber/mo"
	"time"
)

type cacheData[K comparable, T any] struct {
	Mangas map[K]T `json:"mangas"`
}

type cacher[K comparable, T any] struct {
	internal   *cache.Cache[*cacheData[K, T]]
	keyWrapper func(K) K
}

func (c *cacher[K, T]) Get(key K) mo.Option[T] {
	data := c.internal.Get()
	if data.IsPresent() {
		mangas, ok := data.MustGet().Mangas[c.keyWrapper(key)]
		if ok {
			return mo.Some(mangas)
		}
	}

	return mo.None[T]()
}

func (c *cacher[K, T]) Set(key K, t T) error {
	data := c.internal.Get()
	if data.IsPresent() {
		internal := data.MustGet()
		internal.Mangas[c.keyWrapper(key)] = t
		return c.internal.Set(internal)
	} else {
		internal := &cacheData[K, T]{Mangas: make(map[K]T)}
		internal.Mangas[c.keyWrapper(key)] = t
		return c.internal.Set(internal)
	}
}

func (c *cacher[K, T]) Delete(key K) error {
	data := c.internal.Get()
	if data.IsPresent() {
		internal := data.MustGet()
		delete(internal.Mangas, c.keyWrapper(key))
		return c.internal.Set(internal)
	}

	return nil
}

var relationCacher = &cacher[string, int]{
	internal: cache.New(
		"anilist_relation_cache",
		&cache.Options[*cacheData[string, int]]{
			// never expire
			ExpireEvery: constant.Forever,
		},
	),
	keyWrapper: normalizedName,
}

var searchCacher = &cacher[string, []int]{
	internal: cache.New(
		"anilist_search_cache",
		&cache.Options[*cacheData[string, []int]]{
			// update ids every 10 days, since new manga are not added that often
			ExpireEvery: time.Hour * 24 * 10,
		},
	),
	keyWrapper: normalizedName,
}

var idCacher = &cacher[int, *Manga]{
	internal: cache.New(
		"anilist_id_cache",
		&cache.Options[*cacheData[int, *Manga]]{
			// update manga data every day since it can change often
			ExpireEvery: time.Hour * 24,
		},
	),
	keyWrapper: func(id int) int { return id },
}
