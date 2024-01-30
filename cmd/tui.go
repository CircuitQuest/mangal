package cmd

import (
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/tui"
	"github.com/luevano/mangal/tui/state/providers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tuiCmd)
	setDefaultModeShort(tuiCmd)
}

var tuiCmd = &cobra.Command{
	Use:     config.ModeTUI.String(),
	Short:   "Start TUI",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		loaders, err := manager.Loaders(loader.DefaultOptions())
		if err != nil {
			errorf(cmd, err.Error())
		}

		if err := tui.Run(providers.New(loaders)); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
