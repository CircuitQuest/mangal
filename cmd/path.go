package cmd

import (
	"encoding/json"

	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/tui/misc/pathtable"
	"github.com/spf13/cobra"
)

var pathArgs = struct {
	Config    bool
	Cache     bool
	Temp      bool
	Downloads bool
	Providers bool
	Logs      bool
	JSON      bool
}{}

func init() {
	subcommands = append(subcommands, pathCmd)

	pathCmd.Flags().BoolVar(&pathArgs.Config, "config", false, "Path to the config directory")
	pathCmd.Flags().BoolVar(&pathArgs.Cache, "cache", false, "Path to the cache directory")
	pathCmd.Flags().BoolVar(&pathArgs.Temp, "temp", false, "Path to a temporary directory")
	pathCmd.Flags().BoolVar(&pathArgs.Downloads, "downloads", false, "Path to the downloads directory")
	pathCmd.Flags().BoolVar(&pathArgs.Providers, "providers", false, "Path to the providers directory")
	pathCmd.Flags().BoolVar(&pathArgs.Logs, "logs", false, "Path to the logs directory")
	pathCmd.Flags().BoolVarP(&pathArgs.JSON, "json", "j", false, "Output in JSON format for parsing")

	pathCmd.MarkFlagsMutuallyExclusive(
		"config",
		"cache",
		"temp",
		"downloads",
		"providers",
		"logs",
	)
}

// TODO: refactor
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show paths",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		type pathEntry struct {
			Name string `json:"name"`
			Path string `json:"path"`
		}

		var (
			pathToShow     string
			pathToShowName string
		)

		switch {
		case pathArgs.Config:
			// TODO: get config directory based on passed --config flag on the persisted flags
			pathToShow = path.ConfigDir()
			pathToShowName = "config"
		case pathArgs.Providers:
			pathToShow = path.ProvidersDir()
			pathToShowName = "providers"
		case pathArgs.Downloads:
			pathToShow = path.DownloadsDir()
			pathToShowName = "downloads"
		case pathArgs.Cache:
			pathToShow = path.CacheDir()
			pathToShowName = "cache"
		case pathArgs.Temp:
			pathToShow = path.TempDir()
			pathToShowName = "temp"
		case pathArgs.Logs:
			pathToShow = path.LogDir()
			pathToShowName = "logs"
		default:
			if pathArgs.JSON {
				err := json.NewEncoder(cmd.OutOrStdout()).Encode([]pathEntry{
					{
						Name: "config",
						Path: path.ConfigDir(),
					},
					{
						Name: "providers",
						Path: path.ProvidersDir(),
					},
					{
						Name: "downloads",
						Path: path.DownloadsDir(),
					},
					{
						Name: "cache",
						Path: path.CacheDir(),
					},
					{
						Name: "logs",
						Path: path.LogDir(),
					},
					{
						Name: "temp",
						Path: path.TempDir(),
					},
				})
				if err != nil {
					errorf(cmd, err.Error())
				}

				return
			}

			if err := pathtable.Run(); err != nil {
				errorf(cmd, err.Error())
			}

			return
		}

		if pathArgs.JSON {
			err := json.NewEncoder(cmd.OutOrStdout()).Encode(pathEntry{
				Name: pathToShowName,
				Path: pathToShow,
			})
			if err != nil {
				errorf(cmd, err.Error())
			}

			return
		}

		cmd.Println(pathToShow)
	},
}
