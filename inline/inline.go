package inline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/luevano/libmangal"
)

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
	Download bool
	JSON     bool
}

type Options struct {
	InlineArgs
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}

func Run(ctx context.Context, options Options) error {
	queryResult := QueryResult{Query: options.Query, Provider: options.Provider}

	switch {
	case options.Download:
		return fmt.Errorf("unimplemented")
	case options.JSON:
		mangas, err := options.Client.SearchMangas(ctx, options.Query)
		if err != nil {
			return err
		}
		if len(mangas) == 0 {
			return fmt.Errorf("no mangas found with provider ID %q and query %q", options.Provider, options.Query)
		}

		mangaResults := []MangaResult{}
		for i, manga := range mangas {
			anilists, err := options.Anilist.SearchMangas(ctx, manga.Info().AnilistSearch)
			if err != nil {
				return err
			}
			mangaResult := MangaResult{Index: i, Manga: manga}
			if len(anilists) != 0 {
				mangaResult.Anilist = anilists[0]
			}

			mangaResults = append(mangaResults, mangaResult)
		}
		queryResult.Results = mangaResults

		queryResultJSON, err := json.Marshal(queryResult)
		if err != nil {
			return err
		}
		fmt.Println(string(queryResultJSON))
	}

	return nil
}
