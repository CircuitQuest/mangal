package formats

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

type forWhat string

const (
	forRead     forWhat = "read"
	forDownload forWhat = "download"
	forBoth     forWhat = "both"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list   *list.State
	keyMap keyMap
}

// Intermediate implements base.State.
func (*state) Intermediate() bool {
	return true
}

// Backable implements base.State.
func (*state) Backable() bool {
	return true
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (*state) Title() base.Title {
	return base.Title{Text: "Formats"}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *state) Status() string {
	return s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (*state) Init(ctx context.Context) tea.Cmd {
	return nil
}

// Update implements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		i, ok := s.list.SelectedItem().(*item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.setRead):
			return setFormatForCmd(forRead, i.format)
		case key.Matches(msg, s.keyMap.setDownload):
			return setFormatForCmd(forDownload, i.format)
		case key.Matches(msg, s.keyMap.setBoth):
			return setFormatForCmd(forBoth, i.format)
		}
	case formatsUpdatedMsg:
		// TODO: don't write? and only keep the option
		// for the duration of the program?
		return writeConfigCmd
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
