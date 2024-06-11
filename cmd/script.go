package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/script"
	"github.com/luevano/mangal/script/lib"
	"github.com/luevano/mangal/util/afs"
	"github.com/spf13/cobra"
)

var scriptArgs = script.Args{}

func init() {
	rootCmd.AddCommand(scriptCmd)
	f := scriptCmd.Flags()

	f.StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	f.StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	f.BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")
	f.StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")
	f.StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")

	// Reused loader options from inlineCmd
	inlineFlags := inlineCmd.Flags()
	f.AddFlag(inlineFlags.Lookup("nsfw"))
	f.AddFlag(inlineFlags.Lookup("language"))
	f.AddFlag(inlineFlags.Lookup("mangaplus-quality"))
	f.AddFlag(inlineFlags.Lookup("mangadex-data-saver"))
	f.AddFlag(inlineFlags.Lookup("title-chapter-number"))
	f.AddFlag(inlineFlags.Lookup("avoid-duplicate-chapters"))
	f.AddFlag(inlineFlags.Lookup("show-unavailable-chapters"))
	f.AddFlag(inlineFlags.Lookup("parallelism"))
	f.AddFlag(inlineFlags.Lookup("headless-use-flaresolverr"))
	f.AddFlag(inlineFlags.Lookup("headless-flaresolverr-url"))

	scriptCmd.MarkPersistentFlagRequired("provider")
	scriptCmd.MarkPersistentFlagRequired("vars")
	scriptCmd.MarkFlagsOneRequired("file", "string", "stdin")
	scriptCmd.MarkFlagsMutuallyExclusive("file", "string", "stdin")
	scriptCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var scriptCmd = &cobra.Command{
	Use:     config.ModeScript.String(),
	Short:   "Useful for custom process with Lua",
	Long:    fmt.Sprintf("%s, useful for custom process with Lua", config.ModeScript),
	GroupID: groupMode,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		var reader io.Reader

		switch {
		case cmd.Flag("file").Changed:
			file, err := afs.Afero.OpenFile(
				scriptArgs.File,
				os.O_RDONLY,
				config.Download.ModeFile.Get(),
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
	Run: func(cmd *cobra.Command, _ []string) {
		// TODO: use Sprintf?
		filename := fmt.Sprint(meta.AppName, ".lua")
		err := afs.Afero.WriteFile(filename, []byte(lib.LuaDoc()), config.Download.ModeFile.Get())
		if err != nil {
			errorf(cmd, "Error writting library specs: %s", err.Error())
		}
		successf(cmd, "Library specs written to %s", filename)
	},
}
