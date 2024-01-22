package loader

import "github.com/luevano/mangal/config"

type Options struct {
	NSFW                    bool   `json:"nsfw"`
	Language                string `json:"language"`
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
		NSFW:                    config.Config.Filter.NSFW.Get(),
		Language:                config.Config.Filter.Language.Get(),
		MangaDexDataSaver:       config.Config.Filter.TitleChapterNumber.Get(),
		TitleChapterNumber:      config.Config.Filter.TitleChapterNumber.Get(),
		AvoidDuplicateChapters:  config.Config.Filter.AvoidDuplicateChapters.Get(),
		ShowUnavailableChapters: config.Config.Filter.ShowUnavailableChapters.Get(),
		Parallelism:             config.Config.Providers.Parallelism.Get(),
		HeadlessUseFlaresolverr: config.Config.Providers.Headless.UseFlaresolverr.Get(),
		HeadlessFlaresolverrURL: config.Config.Providers.Headless.FlaresolverrURL.Get(),
	}
}
