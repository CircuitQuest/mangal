package viewport

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/ui/color"
)

var _ base.State = (*State)(nil)

type State struct {
	size     base.Size
	title    string
	content  string
	viewport viewport.Model
	padding  base.Size
	keyMap   KeyMap
	styles   Styles
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) Backable() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Viewport", Background: color.Viewport}
}

func (s *State) Subtitle() string {
	return ""
}

func (s *State) Status() string {
	return ""
}

func (s *State) Resize(_size base.Size) {
	s.size = _size
	size := s.paddedSize()

	s.viewport.Width = size.Width
	s.viewport.Height = size.Height
}

func (s *State) Update(model base.Model, msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.viewport, cmd = s.viewport.Update(msg)
	return cmd
}

func (s *State) View(model base.Model) string {
	viewport := fmt.Sprintf("%s\n%s\n%s", s.headerView(), s.viewport.View(), s.footerView())
	return s.styles.ContentWrapper(s.padding.Height, s.padding.Width).Render(viewport)
}

func (s *State) Init(model base.Model) tea.Cmd {
	size := s.paddedSize()

	s.viewport = viewport.New(size.Width, size.Height)
	s.viewport.SetContent(s.styles.Content(size.Width).Render(s.content))
	s.viewport.Style = s.styles.Viewport

	return nil
}

func (s *State) paddedSize() base.Size {
	headerHeight := lipgloss.Height(s.headerView())
	footerHeight := lipgloss.Height(s.footerView())
	verticalMarginHeight := headerHeight + footerHeight

	return base.Size{
		Width:  s.size.Width - (2 * s.padding.Width),
		Height: s.size.Height - (2 * s.padding.Height) - verticalMarginHeight,
	}
}

func (s *State) headerView() string {
	title := s.styles.Title.Render(s.title)
	line := s.styles.Line.Render(fmt.Sprintf("%s%s", strings.Repeat("─", max(0, s.viewport.Width-lipgloss.Width(title)-1)), "╮"))
	space := s.styles.Line.Render(fmt.Sprintf("%s%s", strings.Repeat(" ", max(0, s.viewport.Width-lipgloss.Width(title)-1)), "│"))
	return lipgloss.JoinHorizontal(lipgloss.Bottom, title, lipgloss.JoinVertical(lipgloss.Center, line, space))
}

func (s *State) footerView() string {
	info := s.styles.Info.Render(fmt.Sprintf("%3.f%%", s.viewport.ScrollPercent()*100))
	line := s.styles.Line.Render(fmt.Sprintf("%s%s", "╰", strings.Repeat("─", max(0, s.viewport.Width-lipgloss.Width(info)-1))))
	space := s.styles.Line.Render(fmt.Sprintf("%s%s", "│", strings.Repeat(" ", max(0, s.viewport.Width-lipgloss.Width(info)-1))))
	return lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Center, space, line), info)
}
