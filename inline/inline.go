package inline

import (
	"github.com/luevano/libmangal"
)

var queryResult QueryResult = QueryResult{}

type QueryResult struct {
	QueryParams InlineArgs    `json:"query_params"`
	Results     []MangaResult `json:"results"`
}

type MangaResult struct {
	Index   int                    `json:"index"`
	Manga   libmangal.Manga        `json:"manga"`
	// as a pointer to detect when empty
	Anilist *libmangal.AnilistManga `json:"anilist"`
}

type InlineArgs struct {
	Query          string `json:"query"`
	Provider       string `json:"provider"`
	MangaSelector  string `json:"manga_selector"`
	AnilistID      int    `json:"anilist_id"`
	AnilistDisable bool   `json:"anilist_disable"`
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}
