package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/errorstate"
	"github.com/luevano/mangal/tui/state/providers"
	"github.com/luevano/mangal/tui/state/wrapper/viewport"
)

func Run() error {
	loaders, err := manager.Loaders()
	if err != nil {
		return err
	}
	model := base.New(providers.New(loaders), errorstate.New, viewport.New)
	_, err = tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}
