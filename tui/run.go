package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/program"
	"github.com/luevano/mangal/tui/state/home"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

func Run() error {
	errState := func(err error) base.State {
		title := base.Title{
			Text:       "Error",
			Background: color.Error,
		}
		return viewport.New(title, err.Error(), color.Error)
	}
	logState := func() base.State {
		title := base.Title{
			Text:       "Logs",
			Background: color.Viewport,
		}
		return viewport.New(title, log.Aggregate.String(), color.Viewport)
	}
	model := base.New(home.New(), errState, logState)
	program.SetTUI(tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()))
	_, err := program.TUI().Run()
	return err
}
