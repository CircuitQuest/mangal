package inline

import (
	"context"
	"fmt"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/client/anilist"
	"github.com/luevano/mangal/config"
)

func RunDownload(ctx context.Context, args Args) error {
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
	if len(mangaResults) != 1 {
		return fmt.Errorf("invalid manga selector %q, needs to select 1 manga only", args.MangaSelector)
	}

	if args.AnilistID != 0 {
		err := anilist.Anilist.BindTitleWithID(mangaResults[0].Manga.Info().AnilistSearch, args.AnilistID)
		if err != nil {
			return err
		}
	}

	if err := populateChapters(ctx, client, args, &mangaResults); err != nil {
		return err
	}

	formatOption, err := libmangal.FormatString(args.Format)
	if err != nil {
		return err
	}

	// Take the download options from the config and apply necessary changes
	downloadOptions := config.Config.DownloadOptions()
	downloadOptions.Format = formatOption
	downloadOptions.Directory = args.Directory

	for _, manga := range mangaResults {
		for _, chapter := range *manga.Chapters {
			downloadedPath, err := client.DownloadChapter(ctx, chapter, downloadOptions)
			if err != nil {
				return err
			}
			fmt.Println(downloadedPath)
		}
	}

	return nil
}
