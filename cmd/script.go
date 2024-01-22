package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/anilist"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/icon"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/script"
	"github.com/luevano/mangal/script/lib"
	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
)

var scriptArgs = struct {
	File          string
	String        string
	Stdin         bool
	Provider      string
	Variables     map[string]string
	LoaderOptions *loader.Options
}{
	LoaderOptions: &loader.Options{},
}

func init() {
	subcommands = append(subcommands, scriptCmd)

	scriptCmd.Flags().StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	scriptCmd.Flags().StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	scriptCmd.Flags().BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")

	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")

	scriptCmd.Flags().StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")
	scriptCmd.Flags().StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")

	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.NSFW, "nsfw", config.Config.Filter.NSFW.Get(), "Include NSFW content (when supported)")
	scriptCmd.PersistentFlags().StringVar(&scriptArgs.LoaderOptions.Language, "language", config.Config.Filter.Language.Get(), "Manga/Chapter language")
	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.MangaDexDataSaver, "mangadex-data-saver", config.Config.Filter.MangaDexDataSaver.Get(), "If 'data-saver' should be used (mangadex)")
	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.TitleChapterNumber, "title-chapter-number", config.Config.Filter.TitleChapterNumber.Get(), "If 'Chapter #' should always be included")
	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.AvoidDuplicateChapters, "avoid-duplicate-chapters", config.Config.Filter.AvoidDuplicateChapters.Get(), "Only select one chapter when more are found")
	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.ShowUnavailableChapters, "show-unavailable-chapters", config.Config.Filter.ShowUnavailableChapters.Get(), "If chapter is undownloadable, still show it")
	scriptCmd.PersistentFlags().Uint8Var(&scriptArgs.LoaderOptions.Parallelism, "parallelism", config.Config.Providers.Parallelism.Get(), "Parallelism to use for the provider if supported")
	scriptCmd.PersistentFlags().BoolVar(&scriptArgs.LoaderOptions.HeadlessUseFlaresolverr, "headless-use-flaresolverr", config.Config.Providers.Headless.UseFlaresolverr.Get(), "For providers that use headless, if flaresolverr should be used")
	scriptCmd.PersistentFlags().StringVar(&scriptArgs.LoaderOptions.HeadlessFlaresolverrURL, "headless-flaresolverr-url", config.Config.Providers.Headless.FlaresolverrURL.Get(), "URL for the flaresolverr URL")

	scriptCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var scriptCmd = &cobra.Command{
	Use:     "script",
	Short:   "Run mangal in scripting mode",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var reader io.Reader

		switch {
		case cmd.Flag("file").Changed:
			file, err := afs.Afero.OpenFile(
				scriptArgs.File,
				os.O_RDONLY,
				0755,
			)
			if err != nil {
				errorf(cmd, err.Error())
			}

			defer file.Close()

			reader = file
		case cmd.Flag("string").Changed:
			reader = strings.NewReader(scriptArgs.String)
		case cmd.Flag("stdin").Changed:
			reader = os.Stdin
		default:
			errorf(cmd, "either `file`, `string` or `stdin` is required")
		}

		var options script.Options

		options.Variables = scriptArgs.Variables
		options.Anilist = anilist.Client

		if scriptArgs.Provider != "" {
			client, err := client.NewClientByID(context.Background(), scriptArgs.Provider, *scriptArgs.LoaderOptions)
			if err != nil {
				errorf(cmd, err.Error())
			}
			options.Client = client
		}

		if err := script.Run(context.Background(), reader, options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}

func init() {
	scriptCmd.AddCommand(scriptDocCmd)
}

var scriptDocCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation for the `mangal` lua library",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := lib.Lib(lua.NewState(), lib.Options{})

		filename := fmt.Sprint(l.Name, ".lua")

		err := afs.Afero.WriteFile(filename, []byte(l.LuaDoc()), 0o755)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "%s Library specs written to %s\n", icon.Mark, filename)

		return nil
	},
}
