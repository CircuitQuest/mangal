package chapters

import (
	"fmt"

	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/zyedidia/generic/set"
)

// volume can be nil, which represents a list of chapters for a manga with only one chapter
func New(client *libmangal.Client, manga mangadata.Manga, volume mangadata.Volume, chapters []mangadata.Chapter) *state {
	showVolumeNumber := config.TUI.Chapter.ShowVolumeNumber.Get()
	showChapterNumber := config.TUI.Chapter.ShowNumber.Get()
	showGroup := config.TUI.Chapter.ShowGroup.Get()
	showDate := config.TUI.Chapter.ShowDate.Get()

	keyMap := newKeyMap()
	listWrapper := list.New(
		3,
		"chapter", "chapters",
		chapters,
		func(chapter mangadata.Chapter) _list.DefaultItem {
			volNum := fmt.Sprintf(config.TUI.Chapter.VolumeNumberFormat.Get(), chapter.Volume())
			renderedVolNum := style.Bold.Base.Render(volNum)

			chapNum := fmt.Sprintf(config.TUI.Chapter.NumberFormat.Get(), chapter.Info().Number)
			renderedChapNum := style.Bold.Base.Render(chapNum)

			item := &item{
				chapter:               chapter,
				client:                client,
				renderedVolumeNumber:  renderedVolNum,
				renderedChapterNumber: renderedChapNum,
				showVolumeNumber:      &showVolumeNumber,
				showChapterNumber:     &showChapterNumber,
				showGroup:             &showGroup,
				showDate:              &showDate,
			}
			item.updatePaths()
			item.updateDownloadedFormats()
			item.updateReadAvailablePath()

			return item
		},
		keyMap)

	return &state{
		list:              listWrapper,
		chapters:          chapters,
		volume:            volume,
		manga:             manga,
		client:            client,
		selected:          set.NewMapset[*item](),
		keyMap:            keyMap,
		showVolumeNumber:  &showVolumeNumber,
		showChapterNumber: &showChapterNumber,
		showGroup:         &showGroup,
		showDate:          &showDate,
	}
}
