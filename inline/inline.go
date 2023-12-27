package inline

import (
	"github.com/luevano/libmangal"
)

var queryResult QueryResult = QueryResult{}

type QueryResult struct {
	Query    string        `json:"query"`
	Provider string        `json:"provider"`
	Results  []MangaResult `json:"results"`
}

type MangaResult struct {
	Index   int                    `json:"index"`
	Manga   libmangal.Manga        `json:"manga"`
	Anilist libmangal.AnilistManga `json:"anilist"`
}

type InlineArgs struct {
	Query    string
	Provider string
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}
