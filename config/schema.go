package config

import (
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/icon"
)

type config struct {
	Icons     *registered[string, icon.Type]
	CLI       configCLI
	Read      configRead
	Download  configDownload
	TUI       configTUI
	Providers configProviders
	Library   configLibrary
	Filter    configFilter
}

type configCLI struct {
	ColoredHelp *registered[bool, bool]
	Mode        configCLIMode
}

type configCLIMode struct {
	Default *registered[string, Mode]
}

type configRead struct {
	Format         *registered[string, libmangal.Format]
	History        configReadHistory
	DownloadOnRead *registered[bool, bool]
}

type configReadHistory struct {
	Anilist *registered[bool, bool]
	Local   *registered[bool, bool]
}

type configDownload struct {
	Format       *registered[string, libmangal.Format]
	Path         *registered[string, string]
	Strict       *registered[bool, bool]
	SkipIfExists *registered[bool, bool]
	Manga        configDownloadManga
	Volume       configDownloadVolume
	Chapter      configDownloadChapter
	Metadata     configDownloadMetadata
}

type configDownloadManga struct {
	CreateDir    *registered[bool, bool]
	Cover        *registered[bool, bool]
	Banner       *registered[bool, bool]
	NameTemplate *registered[string, string]
}

type configDownloadVolume struct {
	CreateDir    *registered[bool, bool]
	NameTemplate *registered[string, string]
}

type configDownloadChapter struct {
	NameTemplate *registered[string, string]
}

type configDownloadMetadata struct {
	ComicInfoXML *registered[bool, bool]
	SeriesJSON   *registered[bool, bool]
}

type configTUI struct {
	ExpandSingleVolume *registered[bool, bool]
	Chapter            configTUIChapter
}

type configTUIChapter struct {
	NumberFormat *registered[string, string]
	ShowNumber   *registered[bool, bool]
	ShowDate     *registered[bool, bool]
	ShowGroup    *registered[bool, bool]
}

type configProviders struct {
	Parallelism *registered[int64, uint8]
	Cache       configProvidersCache
	Headless    configProvidersHeadless
}

type configProvidersCache struct {
	TTL *registered[string, string]
}

type configProvidersHeadless struct {
	UseFlaresolverr *registered[bool, bool]
	FlaresolverrURL *registered[string, string]
}

type configLibrary struct {
	Path *registered[string, string]
}

type configFilter struct {
	NSFW                    *registered[bool, bool]
	Language                *registered[string, string]
	MangaDexDataSaver       *registered[bool, bool]
	TitleChapterNumber      *registered[bool, bool]
	AvoidDuplicateChapters  *registered[bool, bool]
	ShowUnavailableChapters *registered[bool, bool]
}
