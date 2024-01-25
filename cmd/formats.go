package cmd

import (
	"github.com/luevano/libmangal"
	"github.com/spf13/cobra"
)

func init() {
	subcommands = append(subcommands, formatsCmd)
}

var formatsCmd = &cobra.Command{
	Use:   "formats",
	Short: "Show available download formats",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		for _, format := range libmangal.FormatStrings() {
			cmd.Println(format)
		}
	},
}
