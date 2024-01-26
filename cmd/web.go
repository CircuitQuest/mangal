package cmd

import (
	"github.com/luevano/mangal/web"
	"github.com/spf13/cobra"
)

func webCmd() *cobra.Command {
	webArgs := web.Args{}

	c := &cobra.Command{
		Use:     "web",
		Short:   "Run mangal with Web UI",
		GroupID: groupMode,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := web.Run(webArgs); err != nil {
				errorf(cmd, err.Error())
			}
		},
	}

	c.Flags().BoolVarP(&webArgs.Open, "open", "o", false, "Open served page in the default browser")
	c.Flags().StringVarP(&webArgs.Port, "port", "p", "6969", "HTTP port to use")

	return c
}
