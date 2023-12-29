package inline

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/luevano/libmangal"
)

func getSelectedMangaResults(mangas []libmangal.Manga, options Options) ([]MangaResult, error) {
	var mangaResults []MangaResult
	switch options.MangaSelector {
	case "all":
		for i, manga := range mangas {
			mangaResults = append(mangaResults, MangaResult{Index: i, Manga: manga})
		}
		return mangaResults, nil
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
			return nil, fmt.Errorf("no manga found with provider %q and exact match %q", options.Provider, options.Query)
		}
		return mangaResults, nil
	default:
		index, err := strconv.Atoi(options.MangaSelector)
		if err != nil {
			return nil, fmt.Errorf("invalid manga selector %q", options.MangaSelector)
		}
		if index < 0 || index >= len(mangas) {
			return nil, fmt.Errorf("invalid manga selector %q (index out of range)", options.MangaSelector)
		}
		mangaResults = []MangaResult{{Index: index, Manga: mangas[index]}}
		return mangaResults, nil
	}
}

func assignAnilist(ctx context.Context, mangaResults *[]MangaResult, options Options) {
	for i, mangaResult := range *mangaResults {
		var anilist libmangal.AnilistManga
		var found bool
		var aniErr error
		if options.AnilistID != -1 {
			anilist, found, aniErr = anilistSearch(ctx, options.Anilist, options.AnilistID)
		} else {
			anilist, found, aniErr = anilistSearch(ctx, options.Anilist, mangaResult.Manga.Info().AnilistSearch)
		}
		if aniErr == nil && found {
			(*mangaResults)[i].Anilist = &anilist
		}
	}
}

// TODO: probably change the return "methodology" (explicit returns)
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

func populateChapters(ctx context.Context, mangaResults *[]MangaResult, options Options) error {
	for i, mangaResult := range *mangaResults {
		mangaVolumesAll, err := options.Client.MangaVolumes(ctx, mangaResult.Manga)
		if err != nil {
			return err
		}
		if len(mangaVolumesAll) == 0 {
			// TODO: use query instead of title?
			return fmt.Errorf("no manga volumes found with provider %q title %q", options.Provider, mangaResult.Manga.Info().Title)
		}

		mangaVolumes, err := getSelectedMangaVolumes(mangaVolumesAll, options)
		if err != nil {
			return err
		}

		volumeResultsAll, err := getVolumeResults(ctx, mangaVolumes, options)
		if err != nil {
			return err
		}

		volumeResults, err := getFilteredVolumes(volumeResultsAll, options)
		if err != nil {
			return err
		}
		(*mangaResults)[i].Volumes = &volumeResults
	}
	return nil
}

func getSelectedMangaVolumes(mangaVolumesAll []libmangal.Volume, options Options) ([]libmangal.Volume, error) {
	switch options.VolumeSelector {
	case "all":
		return mangaVolumesAll, nil
	case "first":
		return []libmangal.Volume{mangaVolumesAll[0]}, nil
	case "last":
		return []libmangal.Volume{mangaVolumesAll[len(mangaVolumesAll)-1]}, nil
	default:
		return nil, fmt.Errorf("invalid volume selector %q", options.VolumeSelector)
	}
}

func getVolumeResults(ctx context.Context, mangaVolumes []libmangal.Volume, options Options) ([]VolumeResult, error) {
	var volumeResults []VolumeResult
	for _, mangaVolume := range mangaVolumes {
		mangaVolumeNumber := mangaVolume.Info().Number
		mangaCAll, err := options.Client.VolumeChapters(ctx, mangaVolume)
		if err != nil {
			return nil, err
		}
		if len(mangaCAll) == 0 {
			return nil, fmt.Errorf("no manga chapters found for volume %d (provider %q, title %q)", mangaVolume.Info().Number, options.Provider, options.Query)
		}

		volumeResults = append(volumeResults, VolumeResult{mangaVolumeNumber, &mangaCAll})
	}
	return volumeResults, nil
}

// could be called getSelectedChapters, but at the end it returns filtered volumes
func getFilteredVolumes(volumeResultsAll []VolumeResult, options Options) ([]VolumeResult, error) {
	switch options.ChapterSelector {
	case "all":
		return volumeResultsAll, nil
	case "first":
		firstV := volumeResultsAll[0]
		firstM := (*firstV.Chapters)[0]
		return []VolumeResult{{firstV.Volume, &[]libmangal.Chapter{firstM}}}, nil
	case "last":
		lastV := volumeResultsAll[len(volumeResultsAll)-1]
		lastM := (*lastV.Chapters)[len(*lastV.Chapters)-1]
		return []VolumeResult{{lastV.Volume, &[]libmangal.Chapter{lastM}}}, nil
	default:
		return nil, fmt.Errorf("invalid chapter selector %q", options.ChapterSelector)
	}
}
