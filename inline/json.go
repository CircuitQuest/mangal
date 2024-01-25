package inline

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/luevano/mangal/client"
)

func RunJSON(ctx context.Context, args Args) error {
	client, err := client.NewClientByID(ctx, args.Provider, *args.LoaderOptions)
	if err != nil {
		return err
	}

	mangas, err := client.SearchMangas(ctx, args.Query)
	if err != nil {
		return err
	}
	if len(mangas) == 0 {
		return fmt.Errorf("no mangas found with provider ID %q and query %q", args.Provider, args.Query)
	}

	mangaResults, err := getSelectedMangaResults(args, mangas)
	if err != nil {
		return err
	}

	if !args.AnilistDisable {
		assignAnilist(ctx, args, &mangaResults)
	}

	if args.ChapterPopulate {
		err := populateChapters(ctx, client, args, &mangaResults)
		if err != nil {
			return err
		}
	}

	queryResult := QueryResult{
		QueryParams: args,
		Results:     mangaResults,
	}
	queryResultJSON, err := json.Marshal(queryResult)
	if err != nil {
		return err
	}
	fmt.Println(string(queryResultJSON))
	return nil
}
