package loading

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mangalorg/mangal/color"
	"github.com/mangalorg/mangal/tui/base"
)

var (
	_ base.State  = (*State)(nil)
	_ help.KeyMap = (*KeyMap)(nil)
)

type KeyMap struct{}

func (k KeyMap) ShortHelp() []key.Binding {
	return nil
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return nil
}

type State struct {
	message string
	spinner spinner.Model
	keyMap  KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Loading", Background: color.Loading}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	//return s.spinner.View()
	return ""
}

func (s *State) Backable() bool {
	return true
}

func (s *State) Resize(size base.Size) {
}

func (s *State) SetMessage(message string) {
	s.message = message
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	s.spinner, cmd = s.spinner.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	return fmt.Sprint(
		lipgloss.NewStyle().Foreground(color.Accent).Render(s.spinner.View()),
		lipgloss.NewStyle().Foreground(color.Secondary).Render(s.message),
	)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return s.spinner.Tick
}
