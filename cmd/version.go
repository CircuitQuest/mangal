package cmd

import (
	"github.com/luevano/mangal/meta"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	versionArgs := struct {
		Short bool
	}{}

	c := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if versionArgs.Short {
				cmd.Println(meta.Version)
				return
			}

			cmd.Println(meta.PrettyVersion())
		},
	}

	c.Flags().BoolVarP(&versionArgs.Short, "short", "s", false, "Only show mangal version number")

	return c
}
