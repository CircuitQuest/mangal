package inline

import (
	"fmt"

	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
)

type QueryResult struct {
	QueryParams Args          `json:"query_params"`
	Results     []MangaResult `json:"results"`
}

// TODO: change Anilist to Metadata in general
type MangaResult struct {
	Index    int                  `json:"index"`
	Manga    mangadata.Manga      `json:"manga"`
	Chapters *[]mangadata.Chapter `json:"chapters"`
	Anilist  *lmanilist.Manga     `json:"anilist"`
}

type Args struct {
	Query                  string `json:"query"`
	Provider               string `json:"provider"`
	MangaSelector          string `json:"manga_selector"`
	ChapterSelector        string `json:"chapter_selector"`
	ChapterPopulate        bool   `json:"chapter_populate"`
	PreferProviderMetadata bool   `json:"prefer_provider_metadata"`
	AnilistID              int    `json:"anilist_id"`
	AnilistDisable         bool   `json:"anilist_disable"`
	JSONOutput             bool   `json:"json_output,omitempty"`
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

func GenericSelectorError(selectorType string, selector string, extraInfo string) string {
	msg := fmt.Sprintf("invalid %s selector %q", selectorType, selector)
	if extraInfo == "" {
		return msg
	} else {
		return fmt.Sprintf("%s (%s)", msg, extraInfo)
	}
}
