package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/program"
	"github.com/luevano/mangal/tui/state/home"
)

func Run() error {
	model := base.New(home.New())
	program.SetTUI(tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()))
	_, err := program.TUI().Run()
	return err
}
