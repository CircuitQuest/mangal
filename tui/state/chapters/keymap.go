package chapters

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

var _ help.KeyMap = (*keyMap)(nil)

func newKeyMap() keyMap {
	return keyMap{
		toggle:              util.Bind("toggle", " "),
		read:                util.Bind("read", "r"),
		download:            util.Bind("download", "d"),
		anilist:             util.Bind("anilist", "A"),
		changeFormat:        util.Bind("change format", "f"),
		openURL:             util.Bind("open url", "o"),
		selectAll:           util.Bind("select all", "a"),
		unselectAll:         util.Bind("unselect all", "backspace"),
		toggleVolumeNumber:  util.Bind("toggle vol num", "v"),
		toggleChapterNumber: util.Bind("toggle number", "c"),
		toggleGroup:         util.Bind("toggle group", "ctrl+g"),
		toggleDate:          util.Bind("toggle date", "ctrl+d"),
	}
}

// keyMap implements help.keyMap.
type keyMap struct {
	toggle,
	read,
	download,
	anilist,
	changeFormat,
	openURL,
	selectAll,
	unselectAll,
	toggleVolumeNumber,
	toggleChapterNumber,
	toggleGroup,
	toggleDate key.Binding
}

// ShortHelp implements help.keyMap.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.toggle,
		k.read,
		k.download,
	}
}

// FullHelp implements help.keyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.ShortHelp(),
		{k.anilist, k.changeFormat, k.openURL},
		{k.selectAll, k.unselectAll},
		{k.toggleVolumeNumber, k.toggleChapterNumber, k.toggleGroup, k.toggleDate},
	}
}
