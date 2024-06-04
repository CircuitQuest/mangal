package inline

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/metadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/client/anilist"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
)

func RunDownload(ctx context.Context, args Args) error {
	// TODO: actually fix the config loading order, this is a temporary hotfix
	args.LoaderOptions.MangaPlusOSVersion = config.Config.Providers.MangaPlus.OSVersion.Get()
	args.LoaderOptions.MangaPlusAppVersion = config.Config.Providers.MangaPlus.AppVersion.Get()
	args.LoaderOptions.MangaPlusAndroidID = config.Config.Providers.MangaPlus.AndroidID.Get()

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
	useMangaMetadata := false
	if args.PreferProviderMetadata {
		if err := manga.Metadata().Validate(); err != nil {
			// TODO: this logger is only really used for TUI, better handle logs
			log.Log("provider metadata is preferred but it's not valid: %s\n", err.Error())
		} else {
			useMangaMetadata = true
		}
	}
	var anilistManga lmanilist.Manga
	var found bool

	// TODO: handle merges/full replacements/fills
	//
	// If AnilistID is provided via argument, then search for it and replace the manga metadata,
	// error out as it is expected to find some metadata (given the id).
	// Otherwise, if prefer provider metadata is true and the metadata is valid, libmangal will search for metadata.
	if !useMangaMetadata && args.AnilistID != 0 {
		anilistManga, found, err = anilist.Anilist.SearchByID(ctx, args.AnilistID)
		if err != nil {
			return err
		}
		if !found {
			return fmt.Errorf("couldn't find anilist manga with id %q (from argument)", args.AnilistID)
		}
		manga.SetMetadata(anilistManga.Metadata())
	}

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
	// SearchMetadata would overwrite provider meta or anilist meta with specified id
	if useMangaMetadata || args.AnilistID != 0 {
		downloadOptions.SearchMetadata = false
	}

	totalChapters := len(chapters)
	retryCount := 0
	for i, chapter := range chapters {
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
			// To avoid abusing the mangaplus api, since there are no status codes returned to check
			// sleep for a second on after each chapter download
			if downChap.ChapterStatus == metadata.DownloadStatusNew && args.Provider == "mango-mangaplus" && i != totalChapters-1 {
				time.Sleep(time.Second)
			}
		}
	}

	return nil
}
