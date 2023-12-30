package inline

import (
	"fmt"

	"github.com/luevano/libmangal"
)

type QueryResult struct {
	QueryParams InlineArgs    `json:"query_params"`
	Results     []MangaResult `json:"results"`
}

type MangaResult struct {
	Index    int                     `json:"index"`
	Manga    libmangal.Manga         `json:"manga"`
	Chapters *[]libmangal.Chapter    `json:"chapters"`
	Anilist  *libmangal.AnilistManga `json:"anilist"`
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

type ChapterSelectorError struct {
	selector  string
	extraInfo string
}

func (sE *ChapterSelectorError) Error() string {
	msg := fmt.Sprintf("invalid chapter selector %q", sE.selector)
	if sE.extraInfo == "" {
		return msg
	} else {
		return fmt.Sprintf("%s (%s)", msg, sE.extraInfo)
	}
}
