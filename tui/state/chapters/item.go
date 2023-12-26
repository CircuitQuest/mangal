package chapters

import (
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/color"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/icon"
	"github.com/luevano/mangal/path"
	"github.com/luevano/libmangal"
	"github.com/zyedidia/generic/set"
)

type Item struct {
	client        *libmangal.Client
	chapter       libmangal.Chapter
	selectedItems *set.Set[*Item]
}

func (i *Item) FilterValue() string {
	return i.chapter.String()
}

func (i *Item) Title() string {
	var title strings.Builder

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
			formatStyle := lipgloss.NewStyle().Bold(true).Foreground(color.Warning)
			title.WriteString(formatStyle.Render(format.String()))
		}
	}

	return title.String()
}

func (i *Item) Description() string {
	return i.chapter.Info().URL
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

func (i *Item) Path(format libmangal.Format) string {
	path := config.Config.Download.Path.Get()

	chapter := i.chapter
	volume := chapter.Volume()
	manga := volume.Manga()

	if config.Config.Download.Manga.CreateDir.Get() {
		path = filepath.Join(path, i.client.ComputeMangaFilename(manga))
	}

	if config.Config.Download.Volume.CreateDir.Get() {
		path = filepath.Join(path, i.client.ComputeVolumeFilename(volume))
	}

	return filepath.Join(path, i.client.ComputeChapterFilename(chapter, format))
}

func (i *Item) IsRecent() bool {
	format := config.Config.Read.Format.Get()
	chapter := i.chapter
	volume := chapter.Volume()
	manga := volume.Manga()

	tempPath := filepath.Join(
		path.TempDir(),
		i.client.ComputeMangaFilename(manga),
		i.client.ComputeVolumeFilename(volume),
		i.client.ComputeChapterFilename(chapter, format),
	)

	exists, err := afs.Afero.Exists(tempPath)
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
