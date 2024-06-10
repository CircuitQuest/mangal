package cmd

import (
	"fmt"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/web"
	"github.com/spf13/cobra"
)

var webArgs = web.Args{}

func init() {
	rootCmd.AddCommand(webCmd)

	f := webCmd.Flags()
	f.BoolVarP(&webArgs.Open, "open", "o", false, "Open served page in the default browser")
	f.StringVarP(&webArgs.Port, "port", "p", "6969", "HTTP port to use")
}

var webCmd = &cobra.Command{
	Use:     config.ModeWeb.String(),
	Short:   "Start server",
	Long:    fmt.Sprintf("%s, start server", config.ModeWeb),
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := web.Run(webArgs); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
