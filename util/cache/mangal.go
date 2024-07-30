package cache

import (
	"github.com/philippgille/gokv"
)

const CacheDBNameMangal = "mangal"

const (
	BucketNameSearchHistory = "search-history"
)

const (
	MangasSearchHistory  = "mangas"
	AnilistSearchHistory = "anilist"
)

var store_ = store{
	openStore: func(bucketName string) (gokv.Store, error) {
		return CacheStore(CacheDBNameMangal, bucketName)
	},
}

type store struct {
	openStore func(bucketName string) (gokv.Store, error)
	store     gokv.Store
}

func (s *store) open(bucketName string) error {
	store, err := s.openStore(bucketName)
	s.store = store
	return err
}

func (s *store) close() error {
	if s.store == nil {
		return nil
	}
	return s.store.Close()
}

// SetMangasSearchHistory will store the manga search history records to the cache.
func SetMangasSearchHistory(records Records) error {
	err := store_.open(BucketNameSearchHistory)
	if err != nil {
		return err
	}
	defer store_.close()

	return store_.store.Set(MangasSearchHistory, records)
}

// GetMangasSearchHistory will populate the given records from the manga search history cache.
func GetMangasSearchHistory(records *Records) (bool, error) {
	err := store_.open(BucketNameSearchHistory)
	if err != nil {
		return false, err
	}
	defer store_.close()

	found, err := store_.store.Get(MangasSearchHistory, records)
	if err != nil {
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}

// SetAnilistSearchHistory will store the manga search history records to the cache.
func SetAnilistSearchHistory(records Records) error {
	err := store_.open(BucketNameSearchHistory)
	if err != nil {
		return err
	}
	defer store_.close()

	return store_.store.Set(AnilistSearchHistory, records)
}

// GetAnilistSearchHistory will populate the given records from the manga search history cache.
func GetAnilistSearchHistory(records *Records) (bool, error) {
	err := store_.open(BucketNameSearchHistory)
	if err != nil {
		return false, err
	}
	defer store_.close()

	found, err := store_.store.Get(AnilistSearchHistory, records)
	if err != nil {
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}
