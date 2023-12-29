package inline

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/luevano/libmangal"
)

func anilistSearch[T string | int](ctx context.Context, aniClient *libmangal.Anilist, queryID T) (aniManga libmangal.AnilistManga, found bool, err error) {
	switch v := reflect.ValueOf(queryID); v.Kind() {
	case reflect.String:
		aniManga, found, err = aniClient.FindClosestManga(ctx, v.String())
		if err != nil {
			return
		}
	case reflect.Int:
		aniManga, found, err = aniClient.GetByID(ctx, int(v.Int()))
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("unexpected error while searching manga on anilist with query/id %q (of type %T)", queryID, queryID)
		return
	}

	if !found {
		err = fmt.Errorf("no manga found on anilist with query/id %q", queryID)
		return
	}
	return
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
	case "all":
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
			return fmt.Errorf("invalid manga selector %q", options.MangaSelector)
		}
		if index < 0 || index >= len(mangas) {
			return fmt.Errorf("invalid manga selector %q (index out of range)", options.MangaSelector)
		}
		mangaResults = []MangaResult{{Index: index, Manga: mangas[index]}}
	}

	if !options.AnilistDisable {
		for i, mangaResult := range mangaResults {
			var anilist libmangal.AnilistManga
			var found bool
			var aniErr error
			if options.AnilistID != -1 {
				anilist, found, aniErr = anilistSearch(ctx, options.Anilist, options.AnilistID)
			} else {
				anilist, found, aniErr = anilistSearch(ctx, options.Anilist, mangaResult.Manga.Info().AnilistSearch)
			}
			if aniErr == nil && found {
				mangaResults[i].Anilist = &anilist
			}
		}
	}

	if options.ChapterPopulate {
		for i, mangaResult := range mangaResults {
			mangaVolumesTemp, err := options.Client.MangaVolumes(ctx, mangaResult.Manga)
			if err != nil {
				return err
			}
			if len(mangaVolumesTemp) == 0 {
				// TODO: use query instead of title?
				return fmt.Errorf("no manga volumes found with provider %q title %q", options.Provider, mangaResult.Manga.Info().Title)
			}

			var mangaVolumes []libmangal.Volume
			switch options.VolumeSelector {
			case "all":
				mangaVolumes = mangaVolumesTemp
			case "first":
				mangaVolumes = []libmangal.Volume{mangaVolumesTemp[0]}
			case "last":
				mangaVolumes = []libmangal.Volume{mangaVolumesTemp[len(mangaVolumesTemp) - 1]}
			default:
				return fmt.Errorf("invalid volume selector %q", options.VolumeSelector)
			}

			var volumeResultsTemp []VolumeResult
			for _, mangaVolume := range mangaVolumes {
				mangaVolumeNumber := mangaVolume.Info().Number
				mangaCAll, err := options.Client.VolumeChapters(ctx, mangaVolume)
				if err != nil {
					return err
				}
				if len(mangaCAll) == 0 {
					// TODO: use query instead of title?
					return fmt.Errorf("no manga chapters found for volume %d (provider %q, title %q)", mangaVolume.Info().Number, options.Provider, mangaResult.Manga.Info().Title)
				}

				volumeResultsTemp = append(volumeResultsTemp, VolumeResult{mangaVolumeNumber, &mangaCAll})
			}

			var volumeResults []VolumeResult
			switch options.ChapterSelector {
			case "all":
				volumeResults = volumeResultsTemp
			case "first":
				firstV := volumeResultsTemp[0]
				firstM := (*firstV.Chapters)[0]
				volumeResults = []VolumeResult{{firstV.Volume, &[]libmangal.Chapter{firstM}}}
			case "last":
				lastV := volumeResultsTemp[len(volumeResultsTemp) - 1]
				lastM := (*lastV.Chapters)[len(*lastV.Chapters) - 1]
				volumeResults = []VolumeResult{{lastV.Volume, &[]libmangal.Chapter{lastM}}}
			default:
				return fmt.Errorf("invalid chapter selector %q", options.ChapterSelector)
			}
			mangaResults[i].Volumes = &volumeResults
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
