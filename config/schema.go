package config

import (
	"io/fs"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/theme/icon"
)

type config struct {
	Icons     *registered[string, icon.Type]
	Cache     configCache
	CLI       configCLI
	Read      configRead
	Download  configDownload
	TUI       configTUI
	Providers configProviders
	Library   configLibrary
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
	UserAgent    *registered[string, string]
	ModeDir      *registered[int64, fs.FileMode]
	ModeFile     *registered[int64, fs.FileMode]
	ModeDB       *registered[int64, fs.FileMode]
	Strict       *registered[bool, bool]
	SkipIfExists *registered[bool, bool]
	Provider     configDownloadProvider
	Manga        configDownloadManga
	Volume       configDownloadVolume
	Chapter      configDownloadChapter
	Metadata     configDownloadMetadata
}

type configDownloadProvider struct {
	CreateDir    *registered[bool, bool]
	NameTemplate *registered[string, string]
}

type configDownloadManga struct {
	CreateDir            *registered[bool, bool]
	Cover                *registered[bool, bool]
	Banner               *registered[bool, bool]
	NameTemplate         *registered[string, string]
	NameTemplateFallback *registered[string, string]
}

type configDownloadVolume struct {
	CreateDir    *registered[bool, bool]
	NameTemplate *registered[string, string]
}

type configDownloadChapter struct {
	NameTemplate *registered[string, string]
}

type configDownloadMetadata struct {
	ComicInfoXML            *registered[bool, bool]
	SeriesJSON              *registered[bool, bool]
	SkipSeriesJSONIfOngoing *registered[bool, bool]
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
	Path        *registered[string, string]
	Parallelism *registered[int64, uint8]
	Headless    configProvidersHeadless
	Filter      configProvidersFilter
	MangaPlus   configProvidersMangaPlus
}

type configCache struct {
	Path *registered[string, string]
	TTL  *registered[string, string]
}

type configProvidersHeadless struct {
	UseFlaresolverr *registered[bool, bool]
	FlaresolverrURL *registered[string, string]
}

type configLibrary struct {
	Path *registered[string, string]
}

type configProvidersFilter struct {
	NSFW                    *registered[bool, bool]
	Language                *registered[string, string]
	MangaPlusQuality        *registered[string, string]
	MangaDexDataSaver       *registered[bool, bool]
	TitleChapterNumber      *registered[bool, bool]
	AvoidDuplicateChapters  *registered[bool, bool]
	ShowUnavailableChapters *registered[bool, bool]
}

type configProvidersMangaPlus struct {
	OSVersion  *registered[string, string]
	AppVersion *registered[string, string]
	AndroidID  *registered[string, string]
}
