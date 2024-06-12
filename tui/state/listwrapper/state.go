package listwrapper

import (
	"fmt"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	stringutil "github.com/luevano/mangal/util/string"
)

var _ base.State = (*State)(nil)

// State implements base.State. Wrapper of list.Model.
type State struct {
	list         list.Model
	notification string
	keyMap       keyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.FilterState() == list.Unfiltered
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "List"}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	singular, plural := s.list.StatusBarItemName()
	subtitle := stringutil.Quantify(len(s.list.VisibleItems()), singular, plural)
	if s.FilterState() == list.FilterApplied {
		return fmt.Sprintf("%s %q", subtitle, s.list.FilterValue())
	}

	return subtitle
}

// Status implements base.State.
func (s *State) Status() string {
	if s.FilterState() == list.Filtering {
		return s.list.FilterInput.View()
	}

	if s.notification != "" {
		return s.list.Paginator.View() + " " + s.notification
	}

	return s.list.Paginator.View()
}

// Resize implements base.State. Wrapper of list.Model.
func (s *State) Resize(size base.Size) {
	s.list.SetSize(size.Width, size.Height)
}

// Init implements base.State.
func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}

// Update implements base.State.
func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case NotificationMsg:
		s.notification = string(msg)
		return nil
	case tea.KeyMsg:
		if s.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, s.keyMap.reverse):
			slices.Reverse(s.Items())
			return tea.Sequence(
				s.list.SetItems(s.Items()),
				s.Notify("Reversed", time.Second),
			)
		}
	}

end:
	s.list, cmd = s.list.Update(msg)
	return cmd
}

// View implements base.State. Wrapper of list.Model.
func (s *State) View(model base.Model) string {
	return s.list.View()
}

// FilterState is a wrapper of list.Model.
func (s *State) FilterState() list.FilterState {
	return s.list.FilterState()
}

// SelectedItem is a wrapper of list.Model.
func (s *State) SelectedItem() list.Item {
	return s.list.SelectedItem()
}

// Items is a wrapper of list.Model.
func (s *State) Items() []list.Item {
	return s.list.Items()
}

// TODO: refactor so that it is usable across all States, not only a list
//
// Notify will show a message in the Status (next to the Title).
func (s *State) Notify(message string, duration time.Duration) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return NotificationMsg(message)
		},
		func() tea.Msg {
			time.Sleep(duration)
			return NotificationMsg("")
		},
	)
}
