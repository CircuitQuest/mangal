package inline

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/luevano/libmangal"
)

func anilistSearch(ctx context.Context, aniClient *libmangal.Anilist, query string) (libmangal.AnilistManga, error) {
	anilists, err := aniClient.SearchMangas(ctx, query)
	if err != nil {
		return libmangal.AnilistManga{}, err
	}
	if len(anilists) == 0 {
		return libmangal.AnilistManga{}, fmt.Errorf("no mangas found on anilist with query %q", query)
	}

	return anilists[0], nil
}

func RunJSON(ctx context.Context, options Options) error {
	queryResult.QueryParams = options.InlineArgs

	mangas, err := options.Client.SearchMangas(ctx, options.Query)
	if err != nil {
		return err
	}
	if len(mangas) == 0 {
		return fmt.Errorf("no mangas found with provider ID %q and query %q", options.Provider, options.Query)
	}

	mangaResults := []MangaResult{}
	switch options.MangaSelector {
	case "", "all":
		for i, manga := range mangas {
			mangaResults = append(mangaResults, MangaResult{Index: i, Manga: manga})
		}
	case "exact":
		ok := false
		for i, manga := range mangas {
			if strings.ToLower(manga.Info().Title) == strings.ToLower(options.Query) {
				mangaResults = []MangaResult{{Index: i, Manga: manga}}
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("no manga found with provider %q and exact match %q", options.Provider, options.Query)
		}
	default:
		index, err := strconv.Atoi(options.MangaSelector)
		if err != nil {
			return fmt.Errorf("invalid manga selector %q (int parse error)", options.MangaSelector)
		}
		if index < 0 {
			return fmt.Errorf("invalid manga selector %q (negative int)", options.MangaSelector)
		}
		if index >= len(mangas) {
			return fmt.Errorf("invalid manga selector %q (index out of range)", options.MangaSelector)
		}

		mangaResults = []MangaResult{{Index: index, Manga: mangas[index]}}
	}

	for i, mangaResult := range mangaResults {
		anilist, err := anilistSearch(ctx, options.Anilist, mangaResult.Manga.Info().AnilistSearch)
		if err == nil {
			mangaResults[i].Anilist = anilist
		}
	}

	queryResult.Results = mangaResults
	queryResultJSON, err := json.Marshal(queryResult)
	if err != nil {
		return err
	}
	fmt.Println(string(queryResultJSON))

	return nil
}
