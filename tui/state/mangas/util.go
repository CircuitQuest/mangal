package mangas

func (s *state) updateItem(item *item) {
	item.updateMetadata()
	item.renderMetadata()
}

func (s *state) updateAllItems() {
	for _, i := range s.list.Items() {
		i := i.(*item)
		s.updateItem(i)
	}
}

func (s *state) updateKeybinds() {
	enable := len(s.list.Items()) != 0
	// enabled based on item availability
	s.keyMap.confirm.SetEnabled(enable)
	s.keyMap.anilist.SetEnabled(enable)
	s.keyMap.metadata.SetEnabled(enable)

	s.keyMap.toggleFullMeta.SetEnabled(enable && *s.extraInfo)
}
