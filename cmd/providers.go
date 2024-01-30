package cmd

import (
	"context"
	"net/url"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/provider/info"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/provider/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(providersCmd)
}

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Providers management",
	Args:  cobra.NoArgs,
}

func init() {
	providersCmd.AddCommand(providersAddCmd)
}

var providersAddCmd = &cobra.Command{
	Use:   "add <url>",
	Short: "Install provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		URL, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		return manager.Add(context.Background(), manager.AddOptions{
			URL: URL,
		})
	},
}

func init() {
	providersCmd.AddCommand(providersUpCmd)
}

var providersUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Update providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Update(context.Background(), manager.UpdateOptions{})
	},
}

func init() {
	providersCmd.AddCommand(providersLsCmd)
}

var providersLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed providers",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		loaders, err := manager.Loaders(loader.DefaultOptions())
		if err != nil {
			return err
		}

		for _, loader := range loaders {
			cmd.Println(loader.Info().ID)
		}

		return nil
	},
}

func init() {
	providersCmd.AddCommand(providersRmCmd)
}

var providersRmCmd = &cobra.Command{
	Use:   "rm tags",
	Short: "Remove provider",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, tag := range args {
			if err := manager.Remove(tag); err != nil {
				errorf(cmd, err.Error())
			}
		}
	},
}

var providersNewArgs = struct {
	Dir string
}{}

func init() {
	providersCmd.AddCommand(providersNewCmd)

	providersNewCmd.Flags().StringVarP(&providersNewArgs.Dir, "dir", "d", path.ProvidersDir(), "directory inside which create a new provider")

	providersNewCmd.MarkFlagDirname("dir")
}

var providersNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new provider",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		options := manager.NewOptions{
			Dir: providersNewArgs.Dir,
			Info: info.Info{
				ProviderInfo: libmangal.ProviderInfo{
					ID:          "test",
					Name:        "test",
					Version:     "0.1.0",
					Description: "Lorem ipsum",
					Website:     "example.com",
				},
				Type: info.TypeLua,
			},
		}

		if err := manager.New(options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
