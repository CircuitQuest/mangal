package chapters

import "github.com/zyedidia/generic/set"

func (s *state) actionRunningNow(action string) {
	s.actionRunning = action
}

func (s *state) updateItem(item *item) {
	item.updateDownloadedFormats()
	item.updateReadAvailablePath()
}

func (s *state) updateItems(items set.Set[*item]) {
	for _, item := range items.Keys() {
		s.updateItem(item)
	}
}
