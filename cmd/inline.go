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

func inlineCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "inline",
		Short:   "Inline mode",
		GroupID: groupMode,
		Args:    cobra.NoArgs,
	}

	// To shorten the statements a bit
	f := c.PersistentFlags()
	lOpts := loader.Options{}

	f.StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	f.StringVarP(&inlineArgs.Provider, "provider", "p", "", "Provider id to use")
	f.StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "all", "Manga selector (all|first|last|exact|<index>)")
	f.StringVarP(&inlineArgs.ChapterSelector, "chapter-selector", "c", "all", "Chapter selector (all|first|last|<index>|[from]-[to])")
	f.IntVarP(&inlineArgs.AnilistID, "anilist-id", "a", 0, "Anilist ID to bind title to")
	setupLoaderOptions(f, &lOpts)
	inlineArgs.LoaderOptions = &lOpts

	c.MarkPersistentFlagRequired("provider")
	c.MarkPersistentFlagRequired("query")
	c.RegisterFlagCompletionFunc("provider", completionProviderIDs)

	// Define inlineJSONCmd outside to be able to mark mutually exclusive flags from parent (this cmd)
	cJSON := inlineJSONCmd()
	c.AddCommand(cJSON)
	cJSON.MarkFlagsMutuallyExclusive("anilist-id", "anilist-disable")

	c.AddCommand(inlineDownloadCmd())

	return c
}

func inlineJSONCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "json",
		Short: "Output search results",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := inline.RunJSON(context.Background(), inlineArgs); err != nil {
				errorf(cmd, err.Error())
			}
		},
	}

	c.Flags().BoolVar(&inlineArgs.ChapterPopulate, "chapter-populate", false, "Populate chapter metadata")
	c.Flags().BoolVar(&inlineArgs.AnilistDisable, "anilist-disable", false, "Disable anilist search")

	return c
}

func inlineDownloadCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "download",
		Short: "Download manga",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := inline.RunDownload(context.Background(), inlineArgs); err != nil {
				errorf(cmd, err.Error())
			}
		},
	}

	formatDesc := fmt.Sprintf("Download format (%s)", strings.Join(libmangal.FormatStrings(), "|"))
	c.Flags().StringVarP(&inlineArgs.Format, "format", "f", config.Config.Download.Format.Get().String(), formatDesc)
	c.Flags().StringVarP(&inlineArgs.Directory, "directory", "d", config.Config.Download.Path.Get(), "Download directory")

	return c
}
