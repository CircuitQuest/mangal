package chapters

import (
	"github.com/luevano/mangal/config"
)

func (s *state) actionRunningNow(action string) {
	s.actionRunning = action
}

func (s *state) updateItem(item *item) {
	item.updatePaths()
	item.updateDownloadedFormats()
	item.updateReadAvailablePath()
}

func (s *state) updateAllItems() {
	for _, i := range s.list.Items() {
		i := i.(*item)
		s.updateItem(i)
	}
}

func (s *state) updateListDelegate() {
	if *s.showDate || *s.showGroup {
		s.list.SetDelegateHeight(3)
	} else {
		s.list.SetDelegateHeight(2)
	}
}

func (s *state) updateRenderedSubtitleFormats() {
	s.renderedSubtitleFormats = s.renderedSep +
		s.styles.subtitle.Render("download ") +
		s.styles.format.Render(config.Download.Format.Get().String()) +
		s.styles.subtitle.Render(" & read ") +
		s.styles.format.Render(config.Read.Format.Get().String())
}

// updateKeybinds enables/disables keybinds whose actions require an item
// (either to perform an action, or change something visually).
func (s *state) updateKeybinds() {
	enable := len(s.list.Items()) != 0
	// require item
	s.keyMap.toggle.SetEnabled(enable)
	s.keyMap.read.SetEnabled(enable)
	s.keyMap.download.SetEnabled(enable)
	s.keyMap.openURL.SetEnabled(enable)
	s.keyMap.selectAll.SetEnabled(enable)
	s.keyMap.unselectAll.SetEnabled(enable)

	// only make sense when there are items
	s.keyMap.toggleVolumeNumber.SetEnabled(enable)
	s.keyMap.toggleChapterNumber.SetEnabled(enable)
	s.keyMap.toggleGroup.SetEnabled(enable)
	s.keyMap.toggleDate.SetEnabled(enable)
}
