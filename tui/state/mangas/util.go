package mangas

import "github.com/charmbracelet/lipgloss"

func (s *state) searchView() string {
	return lipgloss.NewStyle().Padding(0, 1, 1, 1).Render(s.searchInput.View())
}

// updateKeybindings depending on the searching state
func (s *state) updateKeybindings() {
	switch s.searchState {
	case searching, unsearched:
		s.keyMap.confirm.SetEnabled(false)
		s.keyMap.search.SetEnabled(false)
		s.keyMap.cancelSearch.SetEnabled(true)
		s.keyMap.confirmSearch.SetEnabled(s.searchInput.Value() != "")
	default:
		s.keyMap.confirm.SetEnabled(true)
		s.keyMap.search.SetEnabled(true)
		s.keyMap.cancelSearch.SetEnabled(false)
		s.keyMap.confirmSearch.SetEnabled(false)
	}
}
