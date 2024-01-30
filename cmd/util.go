package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/provider/manager"
	"github.com/luevano/mangal/theme/icon"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func successf(cmd *cobra.Command, format string, a ...any) {
	cmd.Printf(fmt.Sprintf("%s %s\n", icon.Check, format), a...)
}

func errorf(cmd *cobra.Command, format string, a ...any) {
	cmd.PrintErrf(fmt.Sprintf("%s %s\n", icon.Cross, format), a...)
	os.Exit(1)
}

func setDefaultModeShort(cmd *cobra.Command) {
	if config.Config.CLI.Mode.Default.Get().String() == cmd.Use {
		cmd.Short = fmt.Sprintf("%s (configured as default)", cmd.Short)
	}
}

func completionProviderIDs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	loaders, err := manager.Loaders(loader.DefaultOptions())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	IDs := lo.Map(loaders, func(loader libmangal.ProviderLoader, _ int) string {
		return loader.Info().ID
	})

	return IDs, cobra.ShellCompDirectiveDefault
}

func completionConfigKeys(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	keys := config.Keys()

	filtered := lo.Filter(keys, func(key string, _ int) bool {
		return strings.HasPrefix(key, toComplete)
	})

	return filtered, cobra.ShellCompDirectiveDefault
}

func setupLoaderOptions(f *pflag.FlagSet, o *loader.Options) {
	c := config.Config.Providers
	f.BoolVar(&o.NSFW, "nsfw", c.Filter.NSFW.Get(), "Include NSFW content (when supported)")
	f.StringVar(&o.Language, "language", c.Filter.Language.Get(), "Manga/Chapter language")
	f.BoolVar(&o.MangaDexDataSaver, "mangadex-data-saver", c.Filter.MangaDexDataSaver.Get(), "Use 'data-saver'")
	f.BoolVar(&o.TitleChapterNumber, "title-chapter-number", c.Filter.TitleChapterNumber.Get(), "Include 'Chapter #' always")
	f.BoolVar(&o.AvoidDuplicateChapters, "avoid-duplicate-chapters", c.Filter.AvoidDuplicateChapters.Get(), "No duplicate chapters")
	f.BoolVar(&o.ShowUnavailableChapters, "show-unavailable-chapters", c.Filter.ShowUnavailableChapters.Get(), "Show undownloadable chapters")
	f.Uint8Var(&o.Parallelism, "parallelism", c.Parallelism.Get(), "Provider parallelism to use (when supported)")
	f.BoolVar(&o.HeadlessUseFlaresolverr, "headless-use-flaresolverr", c.Headless.UseFlaresolverr.Get(), "Use Flaresolverr for headlessproviders")
	f.StringVar(&o.HeadlessFlaresolverrURL, "headless-flaresolverr-url", c.Headless.FlaresolverrURL.Get(), "Flaresolverr service URL")
}
