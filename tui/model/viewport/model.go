package viewport

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	viewport.Model

	title,
	content string
	color lipgloss.Color

	borderHorizontalSize,
	borderVerticalSize int

	style  lipgloss.Style
	keyMap keyMap
}

func (m *Model) Title() string {
	return m.title
}

func (m *Model) Color() lipgloss.Color {
	return m.color
}

func (m *Model) SetData(title, content string, color lipgloss.Color) {
	m.title = title
	m.content = content
	m.Model.SetContent(content)
	m.color = color
	m.style = m.style.BorderForeground(color)
	m.updateKeybinds()
}

func (m *Model) KeyMap() help.KeyMap {
	return m.keyMap
}

func (s *Model) Status() string {
	return fmt.Sprintf("%3.f%%", s.ScrollPercent()*100)
}

func (m *Model) Resize(width, height int) tea.Cmd {
	m.Width = width - m.borderHorizontalSize
	m.Height = height - m.borderVerticalSize
	return nil
}

func (m *Model) Init() tea.Cmd {
	return m.Model.Init()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.copy):
			return func() tea.Msg {
				return clipboard.WriteAll(m.content)
			}
		case key.Matches(msg, m.keyMap.goTop):
			m.GotoTop()
		case key.Matches(msg, m.keyMap.goBottom):
			m.GotoBottom()
		case key.Matches(msg, m.keyMap.back):
			return backCmd
		}
	}
	viewport, cmd := m.Model.Update(msg)
	m.Model = viewport
	return cmd
}

func (m *Model) View() string {
	return m.style.Render(m.Model.View())
}

// updateKeybinds enables/disables keybinds based on the content.
func (m *Model) updateKeybinds() {
	enable := m.content != ""
	m.keyMap.copy.SetEnabled(enable)
	m.keyMap.goTop.SetEnabled(enable)
	m.keyMap.goBottom.SetEnabled(enable)
}
