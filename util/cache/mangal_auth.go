package cache

import (
	"github.com/philippgille/gokv"
	"golang.org/x/oauth2"
)

const CacheDBNameMangalAuth = "mangal-auth"

const (
	BucketNameAuthHistory         = "auth-history"
	BucketNameAnilistAuthData     = "anilist"
	BucketNameMyAnimeListAuthData = "myanimelist"
)

const (
	AnilistAuthHistory     = "anilist"
	MyAnimeListAuthHistory = "myanimelist"
)

var authStore_ = store{
	openStore: func(bucketName string) (gokv.Store, error) {
		return CacheStore(CacheDBNameMangalAuth, bucketName)
	},
}

// SetAuthHistory will store the user history for the provider.
func SetAuthHistory(provider string, userHistory UserHistory) error {
	err := authStore_.open(BucketNameAuthHistory)
	if err != nil {
		return err
	}
	defer authStore_.close()

	return authStore_.store.Set(provider, userHistory)
}

// GetAuthHistory will populate the given user history for the provider.
func GetAuthHistory(provider string, userHistory *UserHistory) (bool, error) {
	err := authStore_.open(BucketNameAuthHistory)
	if err != nil {
		return false, err
	}
	defer authStore_.close()

	found, err := authStore_.store.Get(provider, userHistory)
	if err != nil {
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}

// SetAnilistAuthData will store the oauth2 Token for the username.
func SetAnilistAuthData(username string, token *oauth2.Token) error {
	err := authStore_.open(BucketNameAnilistAuthData)
	if err != nil {
		return err
	}
	defer authStore_.close()

	return authStore_.store.Set(username, *token)
}

// GetAnilistAuthData will populate the given oauth2 Token for the username.
func GetAnilistAuthData(username string, token *oauth2.Token) (bool, error) {
	err := authStore_.open(BucketNameAnilistAuthData)
	if err != nil {
		return false, err
	}
	defer authStore_.close()

	found, err := authStore_.store.Get(username, token)
	if err != nil {
		return false, err
	}
	if found {
		return true, nil
	}
	return false, nil
}

// DeleteAnilistAuthData will delete the oauth2 Token assigned for the username.
//
// Doesn't return an error if not found.
func DeleteAnilistAuthData(username string) error {
	err := authStore_.open(BucketNameAnilistAuthData)
	if err != nil {
		return err
	}
	defer authStore_.close()

	err = authStore_.store.Delete(username)
	if err != nil {
		return err
	}
	return nil
}
