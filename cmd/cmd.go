package cmd

import (
	"fmt"
	"os"
	"strings"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/theme/icon"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const groupMode = "mode"

var rootCmd = &cobra.Command{
	Use:  meta.AppName,
	Args: cobra.NoArgs,
	// A default completion option is always added, this would disable it.
	// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

var subcommands []*cobra.Command

func completionProviderIDs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	loaders, err := manager.Loaders(loader.DefaultOptions())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	IDs := lo.Map(loaders, func(loader libmangal.ProviderLoader, _ int) string {
		return loader.Info().ID
	})

	return IDs, cobra.ShellCompDirectiveDefault
}

func successf(cmd *cobra.Command, format string, a ...any) {
	cmd.Printf(fmt.Sprintf("%s %s\n", icon.Check, format), a...)
}

func errorf(cmd *cobra.Command, format string, a ...any) {
	cmd.PrintErrf(fmt.Sprintf("%s %s\n", icon.Cross, format), a...)
	os.Exit(1)
}

func Execute() {
	var root *cobra.Command

	switch config.Config.CLI.Mode.Default.Get() {
	case config.ModeNone:
		root = rootCmd
	case config.ModeTUI:
		root = tuiCmd
	case config.ModeWeb:
		root = webCmd
	case config.ModeScript:
		root = scriptCmd
	case config.ModeInline:
		root = inlineCmd
	}

	for _, subcommand := range subcommands {
		if subcommand == root {
			continue
		}

		root.AddCommand(subcommand)
	}

	defaultConfiguredMode := ""
	if config.Config.CLI.Mode.Default.Get() != config.ModeNone {
		defaultConfiguredMode = fmt.Sprintf("%s (configured as default)", root.Short)
	}

	root.Use = strings.Replace(root.Use, root.Name(), rootCmd.Name(), 1)
	root.Long = fmt.Sprintf("The ultimate CLI manga downloader\n\n%s", defaultConfiguredMode)
	root.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode Commands:",
	})

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
