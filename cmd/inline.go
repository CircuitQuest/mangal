package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/inline"
	"github.com/luevano/mangal/provider/loader"
	"github.com/spf13/cobra"
)

var inlineArgs = inline.Args{}

func init() {
	subcommands = append(subcommands, inlineCmd)
	setDefaultModeShort(inlineCmd)
	// To shorten the statements a bit
	f := inlineCmd.PersistentFlags()
	lOpts := loader.DefaultOptions()

	f.StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	f.StringVarP(&inlineArgs.Provider, "provider", "p", "", "Provider id to use")
	f.StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "all", "Manga selector (all|first|last|id|exact|closest|<index>)")
	f.StringVarP(&inlineArgs.ChapterSelector, "chapter-selector", "c", "all", "Chapter selector (all|first|last|<num>|[from]-[to])")
	f.IntVarP(&inlineArgs.AnilistID, "anilist-id", "a", 0, "Anilist ID to bind title to")
	setupLoaderOptions(f, &lOpts)
	inlineArgs.LoaderOptions = &lOpts

	inlineCmd.MarkPersistentFlagRequired("provider")
	inlineCmd.MarkPersistentFlagRequired("query")
	inlineCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var inlineCmd = &cobra.Command{
	Use:     config.ModeInline.String(),
	Short:   "Inline, useful for automation",
	GroupID: groupMode,
	Args:    cobra.NoArgs,
}

func init() {
	inlineCmd.AddCommand(inlineJSONCmd)

	inlineJSONCmd.Flags().BoolVar(&inlineArgs.ChapterPopulate, "chapter-populate", false, "Populate chapter metadata")
	inlineJSONCmd.Flags().BoolVar(&inlineArgs.AnilistDisable, "anilist-disable", false, "Disable anilist search")

	inlineJSONCmd.MarkFlagsMutuallyExclusive("anilist-id", "anilist-disable")
}

var inlineJSONCmd = &cobra.Command{
	Use:   "json",
	Short: "Output search results",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := inline.RunJSON(context.Background(), inlineArgs); err != nil {
			errorf(cmd, err.Error())
		}
	},
}

func init() {
	inlineCmd.AddCommand(inlineDownloadCmd)

	formatDesc := fmt.Sprintf("Download format (%s)", strings.Join(libmangal.FormatStrings(), "|"))
	inlineDownloadCmd.Flags().StringVarP(&inlineArgs.Format, "format", "f", "", formatDesc)
	inlineDownloadCmd.Flags().StringVarP(&inlineArgs.Directory, "directory", "d", "", "Download directory")
	inlineDownloadCmd.Flags().BoolVar(&inlineArgs.JSONOutput, "json-output", false, "JSON format for individual chapter download output")
}

var inlineDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download manga",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		if err := inline.RunDownload(context.Background(), inlineArgs); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
