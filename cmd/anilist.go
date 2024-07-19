package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/client/anilist"
	anilistViewer "github.com/luevano/mangal/tui/model/anilist"
	"github.com/spf13/cobra"
)

var anilistLogout bool

func init() {
	rootCmd.AddCommand(anilistCmd)
	anilistCmd.Flags().BoolVarP(&anilistLogout, "logout", "l", false, "Logout of Anilist")
}

var anilistCmd = &cobra.Command{
	Use:   "anilist",
	Short: "Anilist auth commands",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if anilistLogout {
			if err := anilist.Anilist().Logout(); err != nil {
				errorf(cmd, err.Error())
			}

			successf(cmd, "Logged out from Anilist")
			return
		}

		if _, err := tea.NewProgram(anilistViewer.New(anilist.Anilist(), true)).Run(); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
