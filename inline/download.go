package inline

import (
	"context"
	"fmt"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/anilist"
	"github.com/luevano/mangal/client"
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

	// TODO: fix this (should be 0), include it in json.go
	if args.AnilistID != -1 {
		err := anilist.Client.BindTitleWithID(mangaResults[0].Manga.Info().AnilistSearch, args.AnilistID)
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

	// TODO: fix args.Directory using default instead of mangal.toml config set
	downloadOptions := libmangal.DownloadOptions{
		Format:              formatOption,
		Directory:           args.Directory,
		CreateVolumeDir:     config.Config.Download.Volume.CreateDir.Get(),
		CreateMangaDir:      config.Config.Download.Manga.CreateDir.Get(),
		Strict:              config.Config.Download.Strict.Get(),
		SkipIfExists:        config.Config.Download.SkipIfExists.Get(),
		DownloadMangaCover:  config.Config.Download.Manga.Cover.Get(),
		DownloadMangaBanner: config.Config.Download.Manga.Banner.Get(),
		WriteSeriesJson:     config.Config.Download.Metadata.SeriesJSON.Get(),
		WriteComicInfoXml:   config.Config.Download.Metadata.ComicInfoXML.Get(),
		ComicInfoXMLOptions: libmangal.DefaultComicInfoOptions(),
		ImageTransformer: func(bytes []byte) ([]byte, error) {
			return bytes, nil
		},
	}

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
