package loader

import "github.com/luevano/mangal/config"

type Options struct {
	NSFW                    bool   `json:"nsfw"`
	Language                string `json:"language"`
	MangaPlusQuality        string `json:"manga_plus_quality"`
	MangaDexDataSaver       bool   `json:"manga_dex_data_saver"`
	TitleChapterNumber      bool   `json:"title_chapter_number"`
	AvoidDuplicateChapters  bool   `json:"avoid_duplicate_chapters"`
	ShowUnavailableChapters bool   `json:"show_unavailable_chapters"`
	Parallelism             uint8  `json:"parallelism"`
	HeadlessUseFlaresolverr bool   `json:"headless_use_flaresolverr"`
	HeadlessFlaresolverrURL string `json:"headless_flaresolverr_url"`
}

// DefaultOptions gets configured options or default ones.
func DefaultOptions() Options {
	return Options{
		NSFW:                    config.Config.Providers.Filter.NSFW.Get(),
		Language:                config.Config.Providers.Filter.Language.Get(),
		MangaPlusQuality:        config.Config.Providers.Filter.MangaPlusQuality.Get(),
		MangaDexDataSaver:       config.Config.Providers.Filter.MangaDexDataSaver.Get(),
		TitleChapterNumber:      config.Config.Providers.Filter.TitleChapterNumber.Get(),
		AvoidDuplicateChapters:  config.Config.Providers.Filter.AvoidDuplicateChapters.Get(),
		ShowUnavailableChapters: config.Config.Providers.Filter.ShowUnavailableChapters.Get(),
		Parallelism:             config.Config.Providers.Parallelism.Get(),
		HeadlessUseFlaresolverr: config.Config.Providers.Headless.UseFlaresolverr.Get(),
		HeadlessFlaresolverrURL: config.Config.Providers.Headless.FlaresolverrURL.Get(),
	}
}
