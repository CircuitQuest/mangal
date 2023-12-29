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
	Index   int                     `json:"index"`
	Manga   libmangal.Manga         `json:"manga"`
	Volumes *[]VolumeResult         `json:"volumes"`
	Anilist *libmangal.AnilistManga `json:"anilist"`
}

type VolumeResult struct {
	Volume   int                  `json:"volume"`
	Chapters *[]libmangal.Chapter `json:"chapters"`
}

type InlineArgs struct {
	Query           string `json:"query"`
	Provider        string `json:"provider"`
	MangaSelector   string `json:"manga_selector"`
	ChapterSelector string `json:"chapter_selector"`
	ChapterPopulate bool   `json:"chapter_populate"`
	AnilistID       int    `json:"anilist_id"`
	AnilistDisable  bool   `json:"anilist_disable"`
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}
