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
	rootCmd.AddCommand(inlineCmd)
	// To shorten the statements a bit
	f := inlineCmd.PersistentFlags()

	f.StringVarP(&inlineArgs.Query, "query", "q", "", "Query to search")
	f.StringVarP(&inlineArgs.Provider, "provider", "p", "", "Provider id to use")
	f.StringVarP(&inlineArgs.MangaSelector, "manga-selector", "m", "all", "Manga selector (all|first|last|id|exact|closest|<index>)")
	f.StringVarP(&inlineArgs.ChapterSelector, "chapter-selector", "c", "all", "Chapter selector (all|first|last|<num>|[from]-[to])")
	f.BoolVar(&inlineArgs.PreferProviderMetadata, "prefer-provider-metadata", false, "Prefer provider metadata if valid (skips --search-metadata)")
	f.IntVarP(&inlineArgs.AnilistID, "anilist-id", "a", 0, "Anilist ID to search for metadata")
	f.Bool("search-metadata", config.Config.Download.Metadata.Search.Get(), "Search metadata and replace the provider metadata")

	// TODO: remove loaderOptions once script cmd is refactored/removed, these are not used by inline
	loaderOptions := loader.DefaultOptions()
	setupLoaderOptions(f, &loaderOptions)

	inlineCmd.MarkPersistentFlagRequired("provider")
	inlineCmd.MarkPersistentFlagRequired("query")
	inlineCmd.RegisterFlagCompletionFunc("provider", completionProviderIDs)

	// when anilist-id is provided, then use that exclusively
	inlineCmd.MarkFlagsMutuallyExclusive("anilist-id", "search-metadata")
	inlineCmd.MarkFlagsMutuallyExclusive("anilist-id", "prefer-provider-metadata")

	// config(viper) flag binds
	config.BindPFlag(config.Config.Download.Metadata.Search.Key, f.Lookup("search-metadata"))
	config.BindPFlag(config.Config.Providers.Filter.NSFW.Key, f.Lookup("nsfw"))
	config.BindPFlag(config.Config.Providers.Filter.Language.Key, f.Lookup("language"))
	config.BindPFlag(config.Config.Providers.Filter.MangaPlusQuality.Key, f.Lookup("mangaplus-quality"))
	config.BindPFlag(config.Config.Providers.Filter.MangaDexDataSaver.Key, f.Lookup("mangadex-data-saver"))
	config.BindPFlag(config.Config.Providers.Filter.TitleChapterNumber.Key, f.Lookup("title-chapter-number"))
	config.BindPFlag(config.Config.Providers.Filter.AvoidDuplicateChapters.Key, f.Lookup("avoid-duplicate-chapters"))
	config.BindPFlag(config.Config.Providers.Filter.ShowUnavailableChapters.Key, f.Lookup("show-unavailable-chapters"))
	config.BindPFlag(config.Config.Providers.Parallelism.Key, f.Lookup("parallelism"))
	config.BindPFlag(config.Config.Providers.Headless.UseFlaresolverr.Key, f.Lookup("headless-use-flaresolverr"))
	config.BindPFlag(config.Config.Providers.Headless.FlaresolverrURL.Key, f.Lookup("headless-flaresolverr-url"))
}

var inlineCmd = &cobra.Command{
	Use:     config.ModeInline.String(),
	Short:   "Useful for automation",
	Long:    fmt.Sprintf("%s, useful for automation", config.ModeInline),
	GroupID: groupMode,
	Args:    cobra.NoArgs,
}

func init() {
	inlineCmd.AddCommand(inlineJSONCmd)

	f := inlineJSONCmd.Flags()
	f.BoolVar(&inlineArgs.ChapterPopulate, "chapter-populate", false, "Populate chapter metadata")
	f.BoolVar(&inlineArgs.AnilistDisable, "anilist-disable", false, "Disable anilist search")

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

	f := inlineDownloadCmd.Flags()
	fmtDesc := fmt.Sprintf("Download format (%s)", strings.Join(libmangal.FormatStrings(), "|"))
	f.BoolVar(&inlineArgs.JSONOutput, "json-output", false, "JSON format for individual chapter download output")
	f.StringP("format", "f", config.Config.Download.Format.Get().String(), fmtDesc)
	f.StringP("directory", "d", config.Config.Download.Path.Get(), "Download directory")

	inlineDownloadCmd.MarkFlagDirname("directory")

	config.BindPFlag(config.Config.Download.Format.Key, f.Lookup("format"))
	config.BindPFlag(config.Config.Download.Path.Key, f.Lookup("directory"))
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
