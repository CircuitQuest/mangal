package inline

import (
	"context"
	"encoding/json"
	"fmt"
)

func RunJSON(ctx context.Context, options Options) error {
	queryResult.QueryParams = options.InlineArgs

	mangas, err := options.Client.SearchMangas(ctx, options.Query)
	if err != nil {
		return err
	}
	if len(mangas) == 0 {
		return fmt.Errorf("no mangas found with provider ID %q and query %q", options.Provider, options.Query)
	}

	mangaResults, err := getSelectedMangaResults(mangas, options)
	if err != nil {
		return err
	}

	if !options.AnilistDisable {
		assignAnilist(ctx, &mangaResults, options)
	}

	if options.ChapterPopulate {
		populateChapters(ctx, &mangaResults, options)
	}

	queryResult.Results = mangaResults
	queryResultJSON, err := json.Marshal(queryResult)
	if err != nil {
		return err
	}
	fmt.Println(string(queryResultJSON))

	return nil
}
