package cmd

import (
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/provider/loader"
	"github.com/spf13/pflag"
)

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
