package cmd

import (
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/spf13/cobra"
)

const groupMode = "mode"

func init() {
	rootCmd.AddGroup(&cobra.Group{
		ID:    groupMode,
		Title: "Mode Commands:",
	})

	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetIn(os.Stdin)

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
}

var rootCmd = &cobra.Command{
	Use:  meta.AppName,
	Long: "The ultimate CLI manga downloader",
	Args: cobra.NoArgs,
	// A default completion option is always added, this would disable it.
	// CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true}
	// PersistentPreRun: func(cmd *cobra.Command, _ []string){
	// 	cmd.Printf("%s: PersistentPreRun\n", cmd.Name())
	// },
	// Run: func(cmd *cobra.Command, _ []string) {
	// 	// For more on setting a "default cmd":
	// 	// https://github.com/spf13/cobra/issues/823
	// 	// I'm not passing a default command via CLI so it doesn't apply that well
	// 	if config.Config.CLI.Mode.Default.Get() != config.ModeNone {
	// 		cmd.SetArgs([]string{config.Config.CLI.Mode.Default.Get().String()})
	// 		cmd.Execute()
	// 		return
	// 	}
	// 	cmd.Println(config.Config.CLI.Mode.Default.Get())
	// 	cmd.Help()
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorf(rootCmd, err.Error())
	}
}
