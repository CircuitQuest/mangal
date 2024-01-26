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
	pF := c.PersistentFlags()
	cP := config.Config.Providers
	lOpts := loader.Options{}

	pF.StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	pF.StringVarP(&inlineArgs.Provider, "provider", "p", "", "Provider id to use")
	pF.StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "all", "Manga selector (all|first|last|exact|<index>)")
	pF.StringVarP(&inlineArgs.ChapterSelector, "chapter-selector", "c", "all", "Chapter selector (all|first|last|<index>|[from]-[to])")
	pF.IntVarP(&inlineArgs.AnilistID, "anilist-id", "a", 0, "Anilist ID to bind title to")

	// Setup LoaderOptions
	pF.BoolVar(&lOpts.NSFW, "nsfw", cP.Filter.NSFW.Get(), "Include NSFW content (when supported)")
	pF.StringVar(&lOpts.Language, "language", cP.Filter.Language.Get(), "Manga/Chapter language")
	pF.BoolVar(&lOpts.MangaDexDataSaver, "mangadex-data-saver", cP.Filter.MangaDexDataSaver.Get(), "Use 'data-saver'")
	pF.BoolVar(&lOpts.TitleChapterNumber, "title-chapter-number", cP.Filter.TitleChapterNumber.Get(), "Include 'Chapter #' always")
	pF.BoolVar(&lOpts.AvoidDuplicateChapters, "avoid-duplicate-chapters", cP.Filter.AvoidDuplicateChapters.Get(), "No duplicate chapters")
	pF.BoolVar(&lOpts.ShowUnavailableChapters, "show-unavailable-chapters", cP.Filter.ShowUnavailableChapters.Get(), "Show undownloadable chapters")
	pF.Uint8Var(&lOpts.Parallelism, "parallelism", cP.Parallelism.Get(), "Provider parallelism to use (when supported)")
	pF.BoolVar(&lOpts.HeadlessUseFlaresolverr, "headless-use-flaresolverr", cP.Headless.UseFlaresolverr.Get(), "Use Flaresolverr for headlessproviders")
	pF.StringVar(&lOpts.HeadlessFlaresolverrURL, "headless-flaresolverr-url", cP.Headless.FlaresolverrURL.Get(), "Flaresolverr service URL")
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
