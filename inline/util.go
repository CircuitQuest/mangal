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
	case "first":
		return []MangaResult{{Index: 0, Manga: mangas[0]}}, nil
	case "last":
		lastIndex := len(mangas) - 1
		return []MangaResult{{Index: lastIndex, Manga: mangas[lastIndex]}}, nil
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
		volumes, err := options.Client.MangaVolumes(ctx, mangaResult.Manga)
		if err != nil {
			return err
		}
		if len(volumes) == 0 {
			// TODO: use query instead of title?
			return fmt.Errorf("no manga volumes found with provider %q title %q", options.Provider, mangaResult.Manga.Info().Title)
		}

		chapters, err := getChapters(ctx, volumes, options)
		if err != nil {
			return err
		}

		selectedChapters, err := getSelectedChapters(chapters, options)
		if err != nil {
			return err
		}
		(*mangaResults)[i].Chapters = &selectedChapters
	}
	return nil
}

func getChapters(ctx context.Context, volumes []libmangal.Volume, options Options) ([]libmangal.Chapter, error) {
	var chapters []libmangal.Chapter
	for _, volume := range volumes {
		volumeChapters, err := options.Client.VolumeChapters(ctx, volume)
		if err != nil {
			return nil, err
		}
		if len(volumeChapters) == 0 {
			return nil, fmt.Errorf("no manga chapters found for volume %d (provider %q, title %q)", volume.Info().Number, options.Provider, options.Query)
		}

		chapters = append(chapters, volumeChapters...)
	}
	return chapters, nil
}

func getSelectedChapters(chapters []libmangal.Chapter, options Options) ([]libmangal.Chapter, error) {
	switch options.ChapterSelector {
	case "all":
		return chapters, nil
	case "first":
		return []libmangal.Chapter{chapters[0]}, nil
	case "last":
		return []libmangal.Chapter{chapters[len(chapters) - 1]}, nil
	default:
		return nil, fmt.Errorf("invalid chapter selector %q", options.ChapterSelector)
	}
}
