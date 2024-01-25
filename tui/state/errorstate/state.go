package errorstate

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/base"
	"github.com/muesli/reflow/wordwrap"
)

var _ base.State = (*State)(nil)

type State struct {
	error  error
	size   base.Size
	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Error", Background: color.Error}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
	s.size = size
}

func (s *State) Update(model base.Model, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keyMap.Quit):
			return tea.Quit
		case key.Matches(msg, s.keyMap.CopyError):
			return func() tea.Msg {
				return clipboard.WriteAll(s.error.Error())
			}
		}
	}

	return nil
}

func (s *State) View(model base.Model) string {
	wrapped := wordwrap.String(s.error.Error(), int(float64(s.size.Width)/1.2))

	return style.Normal.Error.Render(wrapped)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return nil
}
