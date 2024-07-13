package chapters

import (
	"fmt"

	_list "github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/metadata"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/zyedidia/generic/set"
)

// volume can be nil, which represents a list of chapters for a manga with only one chapter
func New(client *libmangal.Client, manga mangadata.Manga, volume mangadata.Volume, chapters []mangadata.Chapter) *state {
	showVolumeNumber := config.TUI.Chapter.ShowVolumeNumber.Get()
	showChapterNumber := config.TUI.Chapter.ShowNumber.Get()
	showGroup := config.TUI.Chapter.ShowGroup.Get()
	showDate := config.TUI.Chapter.ShowDate.Get()

	_styles := defaultStyles()
	_keyMap := newKeyMap()
	renderedSep := _styles.sep.Render(base.Separator)
	listWrapper := list.New(
		3,
		"chapter", "chapters",
		chapters,
		func(chapter mangadata.Chapter) _list.DefaultItem {
			volNum := fmt.Sprintf(config.TUI.Chapter.VolumeNumberFormat.Get(), chapter.Volume())
			chapNum := fmt.Sprintf(config.TUI.Chapter.NumberFormat.Get(), chapter.Info().Number)

			item := &item{
				chapter:               chapter,
				client:                client,
				renderedSep:           renderedSep,
				renderedVolumeNumber:  volNum,
				renderedChapterNumber: chapNum,
				showVolumeNumber:      &showVolumeNumber,
				showChapterNumber:     &showChapterNumber,
				showGroup:             &showGroup,
				showDate:              &showDate,
				styles:                _styles,
			}
			item.updatePaths()
			item.updateDownloadedFormats()
			item.updateReadAvailablePath()

			return item
		},
		&_keyMap)

	s := &state{
		list:              listWrapper,
		meta:              metadata.New(manga.Metadata()),
		chapters:          chapters,
		volume:            volume,
		manga:             manga,
		client:            client,
		selected:          set.NewMapset[*item](),
		renderedSep:       renderedSep,
		showVolumeNumber:  &showVolumeNumber,
		showChapterNumber: &showChapterNumber,
		showGroup:         &showGroup,
		showDate:          &showDate,
		styles:            _styles,
		keyMap:            &_keyMap,
	}
	s.updateKeybinds()
	return s
}
