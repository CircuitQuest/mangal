package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/tui"
	"github.com/luevano/mangal/tui/state/providers"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const groupMode = "mode"

var rootCmd = &cobra.Command{
	Use:  meta.AppName,
	Long: fmt.Sprintf("%s, CLI manga downloader", meta.AppName),
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		loaders, err := manager.Loaders()
		if err != nil {
			errorf(cmd, err.Error())
		}

		if err := tui.Run(providers.New(loaders)); err != nil {
			errorf(cmd, err.Error())
		}
	},
	// A default completion option is always added, this would disable it.
	// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetIn(os.Stdin)

	rootCmd.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode Commands:",
	})

	// Load the config before any command execution
	// Looks weird, it sets the same variable it uses for default and then path.ConfigDir
	// will also read this set flag but it also creates the directory if doesn't exist.
	rootCmd.PersistentFlags().StringVar(&config.Path, "config", config.Path, "Config file path")
	cobra.OnInitialize(initConfig(rootCmd.PersistentFlags().Lookup("config")))

	if config.Config.CLI.ColoredHelp.Get() {
		cc.Init(&cc.Config{
			RootCmd:         rootCmd,
			Headings:        cc.HiCyan + cc.Bold + cc.Underline,
			Commands:        cc.HiYellow + cc.Bold,
			Example:         cc.Italic,
			ExecName:        cc.Bold,
			Flags:           cc.Bold,
			FlagsDataType:   cc.Italic + cc.HiBlue,
			Aliases:         cc.Italic,
			NoExtraNewlines: true,
			NoBottomNewline: true,
		})
	}
	if err := rootCmd.Execute(); err != nil {
		errorf(rootCmd, err.Error())
	}
}

// re-reads the config from the changed flag
func initConfig(flag *pflag.Flag) func() {
	return func() {
		if flag.Changed {
			if err := config.Load(config.Path); err != nil {
				panic(fmt.Errorf("error loading config from path %s: %s", config.Path, err.Error()))
			}
		}
	}
}
