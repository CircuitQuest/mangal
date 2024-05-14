package chapters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/util/afs"
	"github.com/zyedidia/generic/set"
)

type Item struct {
	client            *libmangal.Client
	chapter           *libmangal.Chapter
	selectedItems     *set.Set[*Item]
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool
	tmpPath           *string
	tmpDownPath       *string
}

func (i *Item) FilterValue() string {
	return (*i.chapter).String()
}

func (i *Item) Title() string {
	var title strings.Builder

	if *i.showChapterNumber {
		chapterNumber := fmt.Sprintf(config.Config.TUI.Chapter.NumberFormat.Get(), (*i.chapter).Info().Number)
		chapterNumberFmt := style.Bold.Base.Render(chapterNumber)
		title.WriteString(chapterNumberFmt)
		title.WriteString(" ")
	}

	title.WriteString(i.FilterValue())

	if i.IsSelected() {
		title.WriteString(" ")
		title.WriteString(icon.Mark.String())
	}

	if i.IsRecent() {
		title.WriteString(" ")
		title.WriteString(icon.Recent.String())
	}

	if formats := i.DownloadedFormats(); formats.Size() > 0 {
		title.WriteString(" ")
		title.WriteString(icon.Download.String())

		// So that formats will be displayed in the same order
		for _, format := range libmangal.FormatValues() {
			if !formats.Has(format) {
				continue
			}

			title.WriteString(" ")
			title.WriteString(style.Bold.Warning.Render(format.String()))
		}
	}

	return title.String()
}

func (i *Item) Description() string {
	var extraInfo strings.Builder

	chapterInfo := (*i.chapter).Info()
	if *i.showDate {
		chapterDate := style.Bold.Secondary.Render(chapterInfo.Date.String())
		extraInfo.WriteString(chapterDate)
		extraInfo.WriteString(" ")
	}

	if *i.showGroup {
		scanlationGroup := style.Italic.Secondary.Render(chapterInfo.ScanlationGroup)
		extraInfo.WriteString(scanlationGroup)
	}

	if extraInfo.String() != "" {
		return fmt.Sprintf(
			"%s\n%s",
			extraInfo.String(),
			chapterInfo.URL)
	}

	return chapterInfo.URL
}

func (i *Item) IsSelected() bool {
	return i.selectedItems.Has(i)
}

func (i *Item) Toggle() {
	if i.IsSelected() {
		i.selectedItems.Remove(i)
	} else {
		i.selectedItems.Put(i)
	}
}

// Updated to avoid computing filenames each frame/update
// TODO: Update tmpPath only when the format changes.
func (i *Item) Path(format libmangal.Format) string {
	// path := config.Config.Download.Path.Get()

	// chapter := i.chapter
	// volume := chapter.Volume()
	// manga := volume.Manga()

	// if config.Config.Download.Provider.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeProviderFilename(i.client.Info()))
	// }

	// if config.Config.Download.Manga.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeMangaFilename(manga))
	// }

	// if config.Config.Download.Volume.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeVolumeFilename(volume))
	// }

	return filepath.Join(*i.tmpDownPath, i.client.ComputeChapterFilename(*i.chapter, format))
}

// Updated to avoid computing filenames each frame/update
// TODO: Update tmpPath only when the format changes.
func (i *Item) IsRecent() bool {
	// format := config.Config.Read.Format.Get()
	// chapter := i.chapter
	// volume := chapter.Volume()
	// manga := volume.Manga()

	// tmpPath := filepath.Join(
	// 	path.TempDir(),
	// 	i.client.ComputeProviderFilename(i.client.Info()),
	// 	i.client.ComputeMangaFilename(manga),
	// 	i.client.ComputeVolumeFilename(volume),
	// 	i.client.ComputeChapterFilename(chapter, format),
	// )

	tmpPath := filepath.Join(*i.tmpPath, i.client.ComputeChapterFilename(*i.chapter, config.Config.Read.Format.Get()))
	exists, err := afs.Afero.Exists(tmpPath)
	if err != nil {
		return false
	}

	return exists
}

func (i *Item) DownloadedFormats() set.Set[libmangal.Format] {
	formats := set.NewMapset[libmangal.Format]()

	for _, format := range libmangal.FormatValues() {
		path := i.Path(format)

		exists, err := afs.Afero.Exists(path)
		if err != nil {
			continue
		}

		if exists {
			formats.Put(format)
		}
	}

	return formats
}
