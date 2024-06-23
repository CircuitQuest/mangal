package chapters

import (
	"github.com/luevano/mangal/config"
	"github.com/zyedidia/generic/set"
)

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
