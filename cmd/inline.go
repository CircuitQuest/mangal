package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/anilist"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/inline"
	"github.com/spf13/cobra"
)

var inlineArgs = inline.InlineArgs{}

func init() {
	subcommands = append(subcommands, inlineCmd)

	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.Provider, "provider", "p", "", "Provider id to use")
	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "all", "Manga selector (all|first|last|exact|<index>)")
	inlineCmd.PersistentFlags().StringVarP(&inlineArgs.ChapterSelector, "chapter-selector", "c", "all", "Chapter selector (all|first|last|<index>|[from]-[to])")
	inlineCmd.PersistentFlags().IntVarP(&inlineArgs.AnilistID, "anilist-id", "a", 0, "Anilist ID to bind title to")

	inlineCmd.MarkPersistentFlagRequired("provider")
	inlineCmd.MarkPersistentFlagRequired("query")

	inlineCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)
}

var inlineCmd = &cobra.Command{
	Use:     "inline",
	Short:   "Inline mode",
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
	Use:     "json",
	Short:   "Output search results in JSON",
	Args:    cobra.NoArgs,
	// TODO: change to RunE
	// TODO: refactor this (similar to Download code)
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClientByID(context.Background(), inlineArgs.Provider)
		if err != nil {
			errorf(cmd, err.Error())
		}

		var options inline.Options

		options.InlineArgs = inlineArgs
		options.Client = client
		options.Anilist = anilist.Client

		if err := inline.RunJSON(context.Background(), options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}

func init() {
	inlineCmd.AddCommand(inlineDownloadCmd)
	formatDesc := fmt.Sprintf("Download format (%s) (defaults to config)", strings.Join(libmangal.FormatStrings(), "|"))
	inlineDownloadCmd.Flags().StringVarP(&inlineArgs.Format, "format", "f", "", formatDesc)
	inlineDownloadCmd.Flags().StringVarP(&inlineArgs.Directory, "directory", "d", "", "Download directory (defaults to config)")
}

var inlineDownloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "Download manga",
	Args:    cobra.NoArgs,
	// TODO: refactor this (similar to JSON code)
	Run: func(cmd *cobra.Command, args []string) {
		client, err := client.NewClientByID(context.Background(), inlineArgs.Provider)
		if err != nil {
			errorf(cmd, err.Error())
		}

		var options inline.Options

		options.InlineArgs = inlineArgs
		options.Client = client
		options.Anilist = anilist.Client

		if err := inline.RunDownload(context.Background(), options); err != nil {
			errorf(cmd, err.Error())
		}
	},
}
