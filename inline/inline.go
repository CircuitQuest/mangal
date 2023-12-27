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
	Anilist libmangal.AnilistManga `json:"anilist"`
}

type InlineArgs struct {
	Query         string `json:"query"`
	Provider      string `json:"provider"`
	MangaSelector string `json:"manga_selector"`
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}
