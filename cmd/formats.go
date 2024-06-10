package cmd

import (
	"github.com/luevano/libmangal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(formatsCmd)
}

var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "Show available download formats",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		for _, format := range libmangal.FormatStrings() {
			cmd.Println(format)
		}
	},
}
