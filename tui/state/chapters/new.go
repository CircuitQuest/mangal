package chapters

import (
	"path/filepath"

	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/util"
	"github.com/zyedidia/generic/set"
)

func New(client *libmangal.Client, manga mangadata.Manga, volume mangadata.Volume, chapters []mangadata.Chapter) *State {
	showChapterNumber := config.TUI.Chapter.ShowNumber.Get()
	showGroup := config.TUI.Chapter.ShowGroup.Get()
	showDate := config.TUI.Chapter.ShowDate.Get()
	selectedSet := set.NewMapset[*Item]()

	listWrapper := list.New(util.NewList(
		3,
		"chapter", "chapters",
		chapters,
		func(chapter mangadata.Chapter) _list.DefaultItem {
			providerFilename := client.ComputeProviderFilename(client.Info())
			mangaFilename := client.ComputeMangaFilename(manga)
			volumeFilename := client.ComputeVolumeFilename(volume)

			tmpPath := filepath.Join(path.TempDir(), providerFilename, mangaFilename, volumeFilename)
			tmpDownPath := path.DownloadsDir()
			if config.Download.Provider.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, providerFilename)
			}
			if config.Download.Manga.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, mangaFilename)
			}
			if config.Download.Volume.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, volumeFilename)
			}

			return &Item{
				chapter:           chapter,
				selectedItems:     &selectedSet,
				client:            client,
				showChapterNumber: &showChapterNumber,
				showGroup:         &showGroup,
				showDate:          &showDate,
				tmpPath:           &tmpPath,
				tmpDownPath:       &tmpDownPath,
			}
		},
	))

	return &State{
		list:     listWrapper,
		chapters: chapters,
		volume:   volume,
		manga:    manga,
		client:   client,
		selected: selectedSet,
		keyMap: keyMap{
			unselectAll:         util.Bind("unselect all", "backspace"),
			selectAll:           util.Bind("select all", "a"),
			toggleChapterNumber: util.Bind("toggle ch num", "c"),
			toggleGroup:         util.Bind("toggle group", "ctrl+g"),
			toggleDate:          util.Bind("toggle date", "ctrl+d"),
			toggle:              util.Bind("toggle", " "),
			read:                util.Bind("read", "r"),
			openURL:             util.Bind("open url", "o"),
			anilist:             util.Bind("anilist", "A"),
			download:            util.Bind("download", "d"),
			confirm:             util.Bind("confirm", "enter"),
			changeFormat:        util.Bind("change format", "f"),
			list:                listWrapper.KeyMap(),
		},
		showChapterNumber: &showChapterNumber,
		showGroup:         &showGroup,
		showDate:          &showDate,
	}
}
