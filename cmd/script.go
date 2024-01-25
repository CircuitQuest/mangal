package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/script"
	"github.com/luevano/mangal/script/lib"
	"github.com/spf13/cobra"
)

var scriptArgs = script.Args{}

func init() {
	subcommands = append(subcommands, scriptCmd)
	// To shorten the statements a bit
	f := scriptCmd.Flags()
	cP := config.Config.Providers
	lOpts := loader.Options{}

	f.StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	f.StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	f.BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")
	f.StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")
	f.StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")

	// Setup LoaderOptions
	f.BoolVar(&lOpts.NSFW, "nsfw", cP.Filter.NSFW.Get(), "Include NSFW content (when supported)")
	f.StringVar(&lOpts.Language, "language", cP.Filter.Language.Get(), "Manga/Chapter language")
	f.BoolVar(&lOpts.MangaDexDataSaver, "mangadex-data-saver", cP.Filter.MangaDexDataSaver.Get(), "Use 'data-saver'")
	f.BoolVar(&lOpts.TitleChapterNumber, "title-chapter-number", cP.Filter.TitleChapterNumber.Get(), "Include 'Chapter #' always")
	f.BoolVar(&lOpts.AvoidDuplicateChapters, "avoid-duplicate-chapters", cP.Filter.AvoidDuplicateChapters.Get(), "No duplicate chapters")
	f.BoolVar(&lOpts.ShowUnavailableChapters, "show-unavailable-chapters", cP.Filter.ShowUnavailableChapters.Get(), "Show undownloadable chapters")
	f.Uint8Var(&lOpts.Parallelism, "parallelism", cP.Parallelism.Get(), "Provider parallelism to use (when supported)")
	f.BoolVar(&lOpts.HeadlessUseFlaresolverr, "headless-use-flaresolverr", cP.Headless.UseFlaresolverr.Get(), "Use Flaresolverr for headlessproviders")
	f.StringVar(&lOpts.HeadlessFlaresolverrURL, "headless-flaresolverr-url", cP.Headless.FlaresolverrURL.Get(), "Flaresolverr service URL")
	scriptArgs.LoaderOptions = &lOpts

	inlineCmd.MarkPersistentFlagRequired("provider")
	inlineCmd.MarkPersistentFlagRequired("vars")
	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")
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
		// TODO: fix mode?
		case cmd.Flag("file").Changed:
			file, err := afs.Afero.OpenFile(
				scriptArgs.File,
				os.O_RDONLY,
				0o755,
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

		if err := script.Run(context.Background(), scriptArgs, reader); err != nil {
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
	Run: func(cmd *cobra.Command, args []string) {
		filename := fmt.Sprint(meta.AppName, ".lua")
		// TODO: fix mode?
		err := afs.Afero.WriteFile(filename, []byte(lib.LuaDoc()), 0o755)
		if err != nil {
			errorf(cmd, "Error writting library specs: %s", err.Error())
		}

		successf(cmd, "Library specs written to %s\n", filename)
	},
}
