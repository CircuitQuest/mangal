package inline

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/luevano/libmangal"
)

func getSelectedMangaResults(mangas []libmangal.Manga, options Options) ([]MangaResult, error) {
	var mangaResults []MangaResult

	totalMangas := len(mangas)
	selector := options.MangaSelector
	switch selector {
	case "all":
		for i, manga := range mangas {
			mangaResults = append(mangaResults, MangaResult{Index: i, Manga: manga})
		}
		return mangaResults, nil
	case "first":
		return []MangaResult{{Index: 0, Manga: mangas[0]}}, nil
	case "last":
		return []MangaResult{{Index: totalMangas - 1, Manga: mangas[totalMangas-1]}}, nil
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
			return nil, &MangaSelectorError{selector, fmt.Sprintf("no manga found with provider %q and exact match %q", options.Provider, options.Query)}
		}
		return mangaResults, nil
	default:
		index, err := strconv.Atoi(selector)
		if err != nil {
			return nil, &MangaSelectorError{selector, err.Error()}
		}
		if index < 0 || index >= totalMangas {
			return nil, &MangaSelectorError{selector, fmt.Sprintf("index out of range(0, %d)", totalMangas-1)}
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
		if options.AnilistID != 0 {
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
	const (
		all   = "all"
		first = "first"
		last  = "last"
		from  = "From"
		to    = "To"
	)

	pattern := fmt.Sprintf(`^(%s|%s|%s|(?P<%s>\d+)?-(?P<%s>\d+)?)$`, all, first, last, from, to)
	selectorRegex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	totalChapters := len(chapters)
	selector := options.ChapterSelector
	switch selector {
	case all:
		return chapters, nil
	case first:
		return chapters[0:1], nil
	case last:
		return chapters[totalChapters-1:], nil
	default:
		if strings.Contains(selector, "-") {
			if !selectorRegex.MatchString(selector) {
				return nil, &ChapterSelectorError{selector, "failed to compile regex, check extra spaces or characters"}
			}
			// default values
			fromI := 0
			toI := totalChapters - 1

			groups := reGroups(selectorRegex, selector)

			fromS := groups[from]
			toS := groups[to]

			if fromS == "" && toS == "" {
				return nil, &ChapterSelectorError{selector, "no 'from' or 'to' specified"}
			}

			if fromS != "" {
				fromITemp, err := strconv.Atoi(fromS)
				if err != nil {
					return nil, &ChapterSelectorError{selector, err.Error()}
				}
				if fromITemp >= totalChapters {
					return nil, &ChapterSelectorError{selector, fmt.Sprintf("'from' index out of range(%d)", totalChapters-1)}
				}
				fromI = fromITemp
			}

			if toS != "" {
				toITemp, err := strconv.Atoi(toS)
				if err != nil {
					return nil, &ChapterSelectorError{selector, err.Error()}
				}
				if toITemp >= totalChapters {
					return nil, &ChapterSelectorError{selector, fmt.Sprintf("'to' index out of range(%d)", totalChapters-1)}
				}
				toI = toITemp
			}

			if fromI > toI {
				return nil, &ChapterSelectorError{selector, "'from' greater than 'to'"}
			}

			return chapters[fromI : toI+1], nil

		} else {
			index, err := strconv.Atoi(selector)
			if err != nil {
				return nil, &ChapterSelectorError{selector, err.Error()}
			}
			if index < 0 || index >= totalChapters {
				return nil, &ChapterSelectorError{selector, fmt.Sprintf("index out of range(0, %d)", totalChapters-1)}
			}
			return []libmangal.Chapter{chapters[index]}, nil
		}
	}
}

func reGroups(pattern *regexp.Regexp, str string) map[string]string {
	groups := make(map[string]string)
	match := pattern.FindStringSubmatch(str)

	for i, name := range pattern.SubexpNames() {
		if i > 0 && i <= len(match) {
			groups[name] = match[i]
		}
	}
	return groups
}
