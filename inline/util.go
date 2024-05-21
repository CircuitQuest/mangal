package inline

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client/anilist"
	"github.com/luevano/mangal/log"
	"github.com/samber/lo"
)

func getSelectedMangaResults(args Args, mangas []libmangal.Manga) ([]MangaResult, error) {
	var mangaResults []MangaResult

	totalMangas := len(mangas)
	selector := strings.ReplaceAll(args.MangaSelector, " ", "")
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
	case "closest":
		var manga *libmangal.Manga
		index := -1
		rank := math.MaxInt

		for i, m := range mangas {
			r := fuzzy.RankMatchNormalizedFold(args.Query, m.Info().Title)
			// Grabs the first ranked match, subsequent
			// with the same rank will not be considered
			if r != -1 && r < rank {
				rank = r
				index = i
				manga = &m
			}
		}
		if index == -1 {
			return nil, &MangaSelectorError{selector, fmt.Sprintf("no manga found with provider %q and closest match %q", args.Provider, args.Query)}
		}

		return []MangaResult{{Index: index, Manga: *manga}}, nil
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

func getChapters(ctx context.Context, client *libmangal.Client, args Args, manga libmangal.Manga) ([]libmangal.Chapter, error) {
	volumes, err := client.MangaVolumes(ctx, manga)
	if err != nil {
		return nil, err
	}
	if len(volumes) == 0 {
		// TODO: use query instead of title?
		return nil, fmt.Errorf("no manga volumes found with provider %q title %q", args.Provider, manga.Info().Title)
	}

	chapters, err := getAllVolumeChapters(ctx, client, args, volumes)
	if err != nil {
		return nil, err
	}

	selectedChapters, err := getSelectedChapters(args, chapters)
	if err != nil {
		return nil, err
	}

	return selectedChapters, nil
}

func getAllVolumeChapters(ctx context.Context, client *libmangal.Client, args Args, volumes []libmangal.Volume) ([]libmangal.Chapter, error) {
	var chapters []libmangal.Chapter
	for _, volume := range volumes {
		volumeChapters, err := client.VolumeChapters(ctx, volume)
		if err != nil {
			return nil, err
		}

		if len(volumeChapters) != 0 {
			chapters = append(chapters, volumeChapters...)
		} else {
			log.Log("No chapters found for volume %.1f", volume.Info().Number)
		}
	}
	return chapters, nil
}

func getSelectedChapters(args Args, chapters []libmangal.Chapter) ([]libmangal.Chapter, error) {
	const (
		selectorAll   = "all"
		selectorFirst = "first"
		selectorLast  = "last"
		selectorFrom  = "from"
		selectorTo    = "to"
	)

	pattern := fmt.Sprintf(`^(%s|%s|%s|%s-%s)$`, selectorAll, selectorFirst, selectorLast, numPattern(selectorFrom), numPattern(selectorTo))
	selectorRegex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Assumes that the chapters are in ascending order
	firstChapNum := chapters[0].Info().Number
	lastChapNum := chapters[len(chapters)-1].Info().Number
	selector := strings.ReplaceAll(args.ChapterSelector, " ", "")
	switch selector {
	case selectorAll:
		return chapters, nil
	case selectorFirst:
		return chapters[0:1], nil
	case selectorLast:
		return chapters[len(chapters)-1:], nil
	default:
		if strings.Contains(selector, "-") {
			if !selectorRegex.MatchString(selector) {
				return nil, &ChapterSelectorError{selector, "failed to match regex pattern to selector"}
			}

			groups := reNamedGroups(selectorRegex, selector)
			fromS := groups[selectorFrom]
			toS := groups[selectorTo]
			if fromS == "" && toS == "" {
				return nil, &ChapterSelectorError{selector, "no 'from' or 'to' specified"}
			}

			from, err := parseNumSelector(selector, "from", fromS, chapters[0].Info().Number, lastChapNum)
			if err != nil {
				return nil, err
			}
			to, err := parseNumSelector(selector, "to", toS, lastChapNum, lastChapNum)
			if err != nil {
				return nil, err
			}

			if from > to {
				return nil, &ChapterSelectorError{selector, fmt.Sprintf("'from' (%s) greater than 'to' (%s)", fmtFloat(from), fmtFloat(to))}
			}

			return lo.Filter(chapters, func(chapter libmangal.Chapter, i int) bool {
				return chapter.Info().Number >= from && chapter.Info().Number <= to
			}), nil

		} else {
			numberTemp, err := strconv.ParseFloat(selector, 32)
			if err != nil {
				return nil, &ChapterSelectorError{selector, err.Error()}
			}
			number := float32(numberTemp)
			if number < firstChapNum || number > lastChapNum {
				return nil, &ChapterSelectorError{selector, fmt.Sprintf("chapter number (%s) out of range(%s, %s)", fmtFloat(number), fmtFloat(firstChapNum), fmtFloat(lastChapNum))}
			}
			// Could return more than one chapter if multiple have the same chapter number for some reason
			return lo.Filter(chapters, func(chapter libmangal.Chapter, i int) bool {
				return chapter.Info().Number == number
			}), nil
		}
	}
}

// try to parse the from/to number selector
func parseNumSelector(selector, matchName, match string, defaultNum, lastChapNum float32) (float32, error) {
	if match != "" {
		fromTemp, err := strconv.ParseFloat(match, 32)
		if err != nil {
			return 0, &ChapterSelectorError{selector, err.Error()}
		}
		if float32(fromTemp) > lastChapNum {
			return 0, &ChapterSelectorError{selector, fmt.Sprintf("'%s' (%s) greater than last chapter number (%s)", matchName, fmtFloat(float32(fromTemp)), fmtFloat(lastChapNum))}
		}
		return float32(fromTemp), nil
	}
	return defaultNum, nil
}

func fmtFloat(n float32) string {
	return strconv.FormatFloat(float64(n), 'f', -1, 32)
}

func numPattern(name string) string {
	return fmt.Sprintf(`(?P<%s>\d+(\.\d+)?)?`, name)
}

func reNamedGroups(pattern *regexp.Regexp, str string) map[string]string {
	groups := make(map[string]string)
	match := pattern.FindStringSubmatch(str)
	for i, value := range match {
		name := pattern.SubexpNames()[i]
		if name != "" {
			groups[name] = value
		}
	}
	return groups
}
