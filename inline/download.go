package inline

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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

	manga := mangaResults[0].Manga
	var anilistManga libmangal.AnilistManga
	var found bool

	if args.AnilistID != 0 {
		anilistManga, found, err = anilist.Anilist.GetByID(ctx, args.AnilistID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("couldn't find anilist manga with id %q", args.AnilistID)
		}
	} else {
		anilistManga, found, err = anilist.Anilist.FindClosestManga(ctx, manga.Info().AnilistSearch)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("couldn't find anilist manga for query %q", manga.Info().AnilistSearch)
		}
	}

	manga.SetAnilistManga(anilistManga)
	chapters, err := getChapters(ctx, client, args, manga)
	if err != nil {
		return err
	}

	// Take the download options from the config and apply necessary changes
	downloadOptions := config.Config.DownloadOptions()
	if args.Format != "" {
		formatOption, err := libmangal.FormatString(args.Format)
		if err != nil {
			return err
		}
		downloadOptions.Format = formatOption
	}
	if args.Directory != "" {
		downloadOptions.Directory = args.Directory
	}

	retryCount := 0
	for _, chapter := range chapters {
		retry := true
		for retry {
			retry = false
			downChap, err := client.DownloadChapter(ctx, chapter, downloadOptions)
			if err != nil {
				errMsg := err.Error()
				// TODO: handle other responses here too if possible
				if strings.Contains(errMsg, "429") && strings.Contains(errMsg, "Retry-After") {
					retry = true
					retryCount += 1
					rcTemp := strings.Split(errMsg, ":")
					retryAfter, err := strconv.Atoi(strings.TrimSpace(rcTemp[len(rcTemp)-1]))
					if err != nil {
						return fmt.Errorf("Error while parsing Retry-Count from error mesage: %s", err.Error())
					}
					fmt.Printf("429 status code while downloading chapter. Waiting %d seconds until retry (retry #%d)\n", retryAfter, retryCount)
					time.Sleep(time.Duration(retryAfter) * time.Second)
					continue
				}
				return err
			}

			if args.JSONOutput {
				dc, err := json.Marshal(downChap)
				if err != nil {
					return err
				}
				fmt.Println(string(dc))
			} else {
				fmt.Println(downChap.Path())
			}
		}
	}

	return nil
}
