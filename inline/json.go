package inline

import (
	"context"
	"encoding/json"
	"fmt"
)

func RunJSON(ctx context.Context, options Options) error {
	queryResult.Query = options.Query
	queryResult.Provider = options.Provider

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

	return nil
}
