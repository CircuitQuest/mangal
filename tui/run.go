package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/tui/model"
	"github.com/luevano/mangal/tui/state/providers"
)

func Run() error {
	loaders, err := manager.Loaders()
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(model.New(providers.New(loaders)), tea.WithAltScreen(), tea.WithMouseCellMotion()).Run()
	return err
}
