package cmd

import (
	"fmt"
	"os"
	"strings"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/path"
	"github.com/spf13/cobra"
)

const groupMode = "mode"

var rootArgs = struct {
	Config string
}{}

var subcommands []*cobra.Command

var rootCmd = &cobra.Command{
	Use:  meta.AppName,
	Args: cobra.NoArgs,
	// A default completion option is always added, this would disable it.
	// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}
}

func Execute() {
	// Actual root Cmd
	var root *cobra.Command

	switch config.Config.CLI.Mode.Default.Get() {
	case config.ModeTUI:
		root = tuiCmd
	case config.ModeWeb:
		root = webCmd
	case config.ModeScript:
		root = scriptCmd
	case config.ModeInline:
		root = inlineCmd
	default:
		// ModeNone basically
		root = rootCmd
	}
	root.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode Commands:",
	})

	// Load the config before any command execution
	root.PersistentFlags().StringVar(&rootArgs.Config, "config", path.ConfigDir(), "Config directory")
	root.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		if err := config.Load(rootArgs.Config); err != nil {
			errorf(cmd, "failed to load config: %s", err.Error())
		}
	}

	for _, subcommand := range subcommands {
		if subcommand == root {
			continue
		}
		root.AddCommand(subcommand)
	}

	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	root.SetIn(os.Stdin)

	if config.Config.CLI.Mode.Default.Get() != config.ModeNone {
		root.Use = strings.Replace(root.Use, root.Name(), rootCmd.Name(), 1)
	}
	root.Long = fmt.Sprintf("The ultimate CLI manga downloader\n\n%s", root.Short)

	if config.Config.CLI.ColoredHelp.Get() {
		cc.Init(&cc.Config{
			RootCmd:         root,
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
	if err := root.Execute(); err != nil {
		errorf(root, err.Error())
	}
}
