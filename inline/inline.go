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

type MangaSelectorError struct {
	selector  string
	extraInfo string
}

func (m *MangaSelectorError) Error() string {
	return GenericSelectorError("manga", m.selector, m.extraInfo)
}

type ChapterSelectorError struct {
	selector  string
	extraInfo string
}

func (m *ChapterSelectorError) Error() string {
	return GenericSelectorError("chapter", m.selector, m.extraInfo)
}

type SelectorError struct {
}

func GenericSelectorError(selectorType string, selector string, extraInfo string) string {
	msg := fmt.Sprintf("invalid %s selector %q", selectorType, selector)
	if extraInfo == "" {
		return msg
	} else {
		return fmt.Sprintf("%s (%s)", msg, extraInfo)
	}
}
