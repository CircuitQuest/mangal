package chapters

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

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
