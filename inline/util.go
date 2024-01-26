package inline

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client/anilist"
)

func getSelectedMangaResults(args Args, mangas []libmangal.Manga) ([]MangaResult, error) {
	var mangaResults []MangaResult

	totalMangas := len(mangas)
	selector := args.MangaSelector
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
			if strings.ToLower(manga.Info().Title) == strings.ToLower(args.Query) {
				mangaResults = []MangaResult{{Index: i, Manga: manga}}
				ok = true
				break
			}
		}
		if !ok {
			return nil, &MangaSelectorError{selector, fmt.Sprintf("no manga found with provider %q and exact match %q", args.Provider, args.Query)}
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

func assignAnilist(ctx context.Context, args Args, mangaResults *[]MangaResult) {
	for i, mangaResult := range *mangaResults {
		var aniManga libmangal.AnilistManga
		var found bool
		var aniErr error
		if args.AnilistID != 0 {
			aniManga, found, aniErr = anilist.Anilist.GetByID(ctx, args.AnilistID)
		} else {
			aniManga, found, aniErr = anilist.Anilist.FindClosestManga(ctx, mangaResult.Manga.Info().AnilistSearch)
		}
		if aniErr == nil && found {
			(*mangaResults)[i].Anilist = &aniManga
		}
	}
}

func populateChapters(ctx context.Context, client *libmangal.Client, args Args, mangaResults *[]MangaResult) error {
	for i, mangaResult := range *mangaResults {
		volumes, err := client.MangaVolumes(ctx, mangaResult.Manga)
		if err != nil {
			return err
		}
		if len(volumes) == 0 {
			// TODO: use query instead of title?
			return fmt.Errorf("no manga volumes found with provider %q title %q", args.Provider, mangaResult.Manga.Info().Title)
		}

		chapters, err := getChapters(ctx, client, args, volumes)
		if err != nil {
			return err
		}

		selectedChapters, err := getSelectedChapters(args, chapters)
		if err != nil {
			return err
		}
		(*mangaResults)[i].Chapters = &selectedChapters
	}
	return nil
}

func getChapters(ctx context.Context, client *libmangal.Client, args Args, volumes []libmangal.Volume) ([]libmangal.Chapter, error) {
	var chapters []libmangal.Chapter
	for _, volume := range volumes {
		volumeChapters, err := client.VolumeChapters(ctx, volume)
		if err != nil {
			return nil, err
		}
		if len(volumeChapters) == 0 {
			return nil, fmt.Errorf("no manga chapters found for volume %d (provider %q, title %q)", volume.Info().Number, args.Provider, args.Query)
		}

		chapters = append(chapters, volumeChapters...)
	}
	return chapters, nil
}

func getSelectedChapters(args Args, chapters []libmangal.Chapter) ([]libmangal.Chapter, error) {
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
	selector := args.ChapterSelector
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
