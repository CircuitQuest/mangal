package chapters

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap(listKeyMap help.KeyMap) keyMap {
	return keyMap{
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
		list:                listKeyMap,
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	unselectAll,
	selectAll,
	toggleChapterNumber,
	toggleGroup,
	toggleDate,
	toggle,
	read,
	openURL,
	download,
	anilist,
	confirm,
	changeFormat key.Binding

	list help.KeyMap
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return append(
		k.list.ShortHelp(),
		k.toggle,
		k.read,
		k.download,
	)
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{k.selectAll, k.unselectAll, k.toggleChapterNumber, k.toggleGroup, k.toggleDate},
		{k.anilist},
		{k.changeFormat, k.openURL},
	}
}
