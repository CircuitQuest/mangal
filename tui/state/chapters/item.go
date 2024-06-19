package chapters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/util/afs"
	"github.com/zyedidia/generic/set"
)

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	chapter           mangadata.Chapter
	client            *libmangal.Client
	selectedItems     *set.Set[*item]
	showVolumeNumber  *bool
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool
	tmpPath           *string
	tmpDownPath       *string
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	return i.chapter.String()
}

// Title implements list.DefaultItem.
func (i *item) Title() string {
	var title strings.Builder

	if *i.showVolumeNumber {
		volumeNumber := fmt.Sprintf(config.TUI.Chapter.VolumeNumberFormat.Get(), i.chapter.Volume())
		volumeNumberFmt := style.Bold.Base.Render(volumeNumber)
		title.WriteString(volumeNumberFmt)
		title.WriteString(" ")
	}

	if *i.showChapterNumber {
		chapterNumber := fmt.Sprintf(config.TUI.Chapter.NumberFormat.Get(), i.chapter.Info().Number)
		chapterNumberFmt := style.Bold.Base.Render(chapterNumber)
		title.WriteString(chapterNumberFmt)
		title.WriteString(" ")
	}

	title.WriteString(i.FilterValue())

	if i.isSelected() {
		title.WriteString(" ")
		title.WriteString(icon.Mark.String())
	}

	if i.isRecent() {
		title.WriteString(" ")
		title.WriteString(icon.Recent.String())
	}

	if formats := i.downloadedFormats(); formats.Size() > 0 {
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

// Description implements list.DefaultItem.
func (i *item) Description() string {
	var extraInfo strings.Builder

	chapterInfo := i.chapter.Info()
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

func (i *item) isSelected() bool {
	return i.selectedItems.Has(i)
}

func (i *item) toggle() {
	if i.isSelected() {
		i.selectedItems.Remove(i)
	} else {
		i.selectedItems.Put(i)
	}
}

// TODO: update tmpPath only when the format changes
func (i *item) path(format libmangal.Format) string {
	// path := config.Download.Path.Get()

	// chapter := i.chapter
	// volume := chapter.Volume()
	// manga := volume.Manga()

	// if config.Download.Provider.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeProviderFilename(i.client.Info()))
	// }

	// if config.Download.Manga.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeMangaFilename(manga))
	// }

	// if config.Download.Volume.CreateDir.Get() {
	// 	path = filepath.Join(path, i.client.ComputeVolumeFilename(volume))
	// }

	return filepath.Join(*i.tmpDownPath, i.client.ComputeChapterFilename(i.chapter, format))
}

// TODO: update tmpPath only when the format changes
func (i *item) isRecent() bool {
	// format := config.Read.Format.Get()
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

	tmpPath := filepath.Join(*i.tmpPath, i.client.ComputeChapterFilename(i.chapter, config.Read.Format.Get()))
	exists, err := afs.Afero.Exists(tmpPath)
	if err != nil {
		return false
	}

	return exists
}

func (i *item) downloadedFormats() set.Set[libmangal.Format] {
	formats := set.NewMapset[libmangal.Format]()

	for _, format := range libmangal.FormatValues() {
		path := i.path(format)

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
