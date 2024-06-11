package chapters

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/util"
	"github.com/zyedidia/generic/set"
)

func New(client *libmangal.Client, manga *mangadata.Manga, volume *mangadata.Volume, chapters []*mangadata.Chapter) *State {
	showChapterNumber := config.TUI.Chapter.ShowNumber.Get()
	showGroup := config.TUI.Chapter.ShowGroup.Get()
	showDate := config.TUI.Chapter.ShowDate.Get()
	selectedSet := set.NewMapset[*Item]()

	listWrapper := listwrapper.New(util.NewList(
		3,
		"chapter", "chapters",
		chapters,
		func(chapter *mangadata.Chapter) list.DefaultItem {
			tmpPath := filepath.Join(
				path.TempDir(),
				client.ComputeProviderFilename(client.Info()),
				client.ComputeMangaFilename(*manga),
				client.ComputeVolumeFilename(*volume),
			)

			tmpDownPath := path.DownloadsDir()
			if config.Download.Provider.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, client.ComputeProviderFilename(client.Info()))
			}
			if config.Download.Manga.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, client.ComputeMangaFilename(*manga))
			}
			if config.Download.Volume.CreateDir.Get() {
				tmpDownPath = filepath.Join(tmpDownPath, client.ComputeVolumeFilename(*volume))
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
		client:   client,
		manga:    manga,
		volume:   volume,
		selected: selectedSet,
		list:     listWrapper,
		keyMap: KeyMap{
			UnselectAll:         util.Bind("unselect all", "backspace"),
			SelectAll:           util.Bind("select all", "a"),
			ToggleChapterNumber: util.Bind("toggle ch num", "c"),
			ToggleGroup:         util.Bind("toggle group", "ctrl+g"),
			ToggleDate:          util.Bind("toggle date", "ctrl+d"),
			Toggle:              util.Bind("toggle", " "),
			Read:                util.Bind("read", "r"),
			OpenURL:             util.Bind("open url", "o"),
			Anilist:             util.Bind("anilist", "A"),
			Download:            util.Bind("download", "d"),
			Confirm:             util.Bind("confirm", "enter"),
			ChangeFormat:        util.Bind("change format", "f"),
			list:                listWrapper.GetKeyMap(),
		},
		showChapterNumber: &showChapterNumber,
		showGroup:         &showGroup,
		showDate:          &showDate,
	}
}
