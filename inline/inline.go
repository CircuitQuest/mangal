package inline

import (
	"fmt"

	"github.com/luevano/libmangal/mangadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/provider/loader"
)

type QueryResult struct {
	QueryParams Args          `json:"query_params"`
	Results     []MangaResult `json:"results"`
}

type MangaResult struct {
	Index    int                  `json:"index"`
	Manga    mangadata.Manga      `json:"manga"`
	Chapters *[]mangadata.Chapter `json:"chapters"`
	Anilist  *lmanilist.Manga       `json:"anilist"`
}

type Args struct {
	Query           string          `json:"query"`
	Provider        string          `json:"provider"`
	MangaSelector   string          `json:"manga_selector"`
	ChapterSelector string          `json:"chapter_selector"`
	ChapterPopulate bool            `json:"chapter_populate"`
	AnilistID       int             `json:"anilist_id"`
	AnilistDisable  bool            `json:"anilist_disable"`
	Format          string          `json:"format,omitempty"`
	Directory       string          `json:"directory,omitempty"`
	JSONOutput      bool            `json:"json_output,omitempty"`
	LoaderOptions   *loader.Options `json:"loader_options"`
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

type SelectorError struct{}

func GenericSelectorError(selectorType string, selector string, extraInfo string) string {
	msg := fmt.Sprintf("invalid %s selector %q", selectorType, selector)
	if extraInfo == "" {
		return msg
	} else {
		return fmt.Sprintf("%s (%s)", msg, extraInfo)
	}
}
