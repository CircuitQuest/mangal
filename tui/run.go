package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/errorstate"
	"github.com/luevano/mangal/tui/state/home"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

func Run() error {
	model := base.New(home.New(), errorstate.New, viewport.New)
	_, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}
