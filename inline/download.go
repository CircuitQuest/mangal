package inline

import (
	"context"
	"fmt"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
)

func RunDownload(ctx context.Context, options Options) error {
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
	if len(mangaResults) != 1 {
		return fmt.Errorf("invalid manga selector %q, needs to select 1 manga only", options.MangaSelector)
	}

	if options.AnilistID != -1 {
		err := options.Anilist.BindTitleWithID(mangaResults[0].Manga.Info().AnilistSearch, options.AnilistID)
		if err != nil {
			return err
		}
	}

	if err := populateChapters(ctx, &mangaResults, options); err != nil {
		return err
	}

	downloadOptions := libmangal.DownloadOptions{
		Format:              config.Config.Download.Format.Get(),
		Directory:           config.Config.Download.Path.Get(),
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
			downloadedPath, err := options.Client.DownloadChapter(ctx, chapter, downloadOptions)
			if err != nil {
				return err
			}
			fmt.Println(downloadedPath)
		}
	}

	return nil
}
