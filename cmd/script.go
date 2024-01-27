package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/script"
	"github.com/luevano/mangal/script/lib"
	"github.com/luevano/mangal/util/afs"
	"github.com/spf13/cobra"
)

func scriptCmd() *cobra.Command {
	scriptArgs := script.Args{}

	c := &cobra.Command{
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
					path.ModeFile,
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
	// To shorten the statements a bit
	f := c.Flags()
	lOpts := loader.Options{}

	f.StringVarP(&scriptArgs.File, "file", "f", "", "Read script from file")
	f.StringVarP(&scriptArgs.String, "string", "s", "", "Read script from script")
	f.BoolVarP(&scriptArgs.Stdin, "stdin", "i", false, "Read script from stdin")
	f.StringVarP(&scriptArgs.Provider, "provider", "p", "", "Load provider by tag")
	f.StringToStringVarP(&scriptArgs.Variables, "vars", "v", nil, "Variables to set in the `Vars` table")
	setupLoaderOptions(f, &lOpts)
	scriptArgs.LoaderOptions = &lOpts

	c.MarkPersistentFlagRequired("provider")
	c.MarkPersistentFlagRequired("vars")
	c.MarkFlagsOneRequired("file", "string", "stdin")
	c.MarkFlagsMutuallyExclusive("file", "string", "stdin")
	c.RegisterFlagCompletionFunc("provider", completionProviderIDs)

	c.AddCommand(scriptDocCmd())

	return c
}

func scriptDocCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doc",
		Short: "Generate documentation for the `mangal` lua library",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			filename := fmt.Sprint(meta.AppName, ".lua")
			err := afs.Afero.WriteFile(filename, []byte(lib.LuaDoc()), path.ModeFile)
			if err != nil {
				errorf(cmd, "Error writting library specs: %s", err.Error())
			}
			successf(cmd, "Library specs written to %s", filename)
		},
	}
}
