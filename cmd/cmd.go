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

type Cmd struct {
	root   *cobra.Command
	tui    *cobra.Command
	web    *cobra.Command
	script *cobra.Command
	inline *cobra.Command

	subcommands []*cobra.Command
}

func NewCmd() Cmd {
	c := Cmd{
		root: &cobra.Command{
			Use:  meta.AppName,
			Args: cobra.NoArgs,
			// A default completion option is always added, this would disable it.
			// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		},
		tui: tuiCmd(),
		web: webCmd(),
		script: scriptCmd(),
		inline: inlineCmd(),
	}
	c.subcommands = append(c.subcommands, c.tui)
	c.subcommands = append(c.subcommands, c.web)
	c.subcommands = append(c.subcommands, c.script)
	c.subcommands = append(c.subcommands, c.inline)

	c.subcommands = append(c.subcommands, versionCmd())
	c.subcommands = append(c.subcommands, templatesCmd())
	c.subcommands = append(c.subcommands, formatsCmd())
	c.subcommands = append(c.subcommands, pathCmd())
	c.subcommands = append(c.subcommands, providersCmd())
	c.subcommands = append(c.subcommands, anilistCmd())
	c.subcommands = append(c.subcommands, configCmd())

	return c
}

func (c *Cmd) AddSubcommand(subcommand *cobra.Command) {
	c.subcommands = append(c.subcommands, subcommand)
}

func (c *Cmd) Execute() {
	// Keep a reference to the starter root
	starterRoot := c.root
	switch config.Config.CLI.Mode.Default.Get() {
	case config.ModeNone:
		// Keep the root cmd as is
	case config.ModeTUI:
		c.root = c.tui
	case config.ModeWeb:
		c.root = c.web
	case config.ModeScript:
		c.root = c.script
	case config.ModeInline:
		c.root = c.inline
	}

	for _, subcommand := range c.subcommands {
		if subcommand == c.root {
			continue
		}

		c.root.AddCommand(subcommand)
	}

	defaultConfiguredMode := ""
	if config.Config.CLI.Mode.Default.Get() != config.ModeNone {
		defaultConfiguredMode = fmt.Sprintf("%s (configured as default)", c.root.Short)
	}

	c.root.Use = strings.Replace(c.root.Use, c.root.Name(), starterRoot.Name(), 1)
	c.root.Long = fmt.Sprintf("The ultimate CLI manga downloader\n\n%s", defaultConfiguredMode)
	c.root.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode Commands:",
	})

	if config.Config.CLI.ColoredHelp.Get() {
		cc.Init(&cc.Config{
			RootCmd:         c.root,
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

	c.root.SetOut(os.Stdout)
	c.root.SetErr(os.Stderr)
	c.root.SetIn(os.Stdin)
	if err := c.root.Execute(); err != nil {
		errorf(c.root, err.Error())
	}
}
