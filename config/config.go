package config

import (
	"io/fs"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"github.com/adrg/xdg"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/template/funcs"
	"github.com/luevano/mangal/theme/icon"
)

// Dir is the default config directory which is also set by the --config flag.
//
// This should only be used inside cmd/cmd.go, anywhere else path.Dir() should be used
// as it uses this variable but also has the option to create the directories if needed.
var Dir string = func() string {
	var dir string
	if runtime.GOOS == "darwin" {
		dir = filepath.Join(xdg.Home, ".config", meta.AppName)
	} else {
		dir = filepath.Join(xdg.ConfigHome, meta.AppName)
	}

	return dir
}()

// Config registers all possible config options and it's exported.
var Config = config{
	Icons: reg(field[string, icon.Type]{
		Key:         "icons",
		Default:     icon.TypeASCII,
		Description: "Icon format to use.",
		Unmarshal: func(s string) (icon.Type, error) {
			return icon.TypeString(s)
		},
		Marshal: func(i icon.Type) (string, error) {
			return i.String(), nil
		},
	}),
	Cache: configCache{
		Path: reg(field[string, string]{
			Key:         "cache.path",
			Default:     filepath.Join(xdg.CacheHome, meta.AppName),
			Description: "Path where cache will be written to.",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
		TTL: reg(field[string, string]{
			Key:         "cache.ttl",
			Default:     "24h",
			Description: `Time to live. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
			Validate: func(s string) error {
				_, err := time.ParseDuration(s)
				return err
			},
		}),
	},
	CLI: configCLI{
		ColoredHelp: reg(field[bool, bool]{
			Key:         "cli.colored_help",
			Default:     true,
			Description: "Enable colors in cli help.",
		}),
		Mode: configCLIMode{
			Default: reg(field[string, Mode]{
				Key:         "cli.mode.default",
				Default:     ModeTUI,
				Description: "Default mode to use when no subcommand is given.",
				Unmarshal: func(s string) (Mode, error) {
					return ModeString(s)
				},
				Marshal: func(mode Mode) (string, error) {
					return mode.String(), nil
				},
			}),
		},
	},
	Read: configRead{
		Format: reg(field[string, libmangal.Format]{
			Key:         "read.format",
			Default:     libmangal.FormatPDF,
			Description: "Format to read chapters in.",
			Unmarshal: func(s string) (libmangal.Format, error) {
				return libmangal.FormatString(s)
			},
			Marshal: func(format libmangal.Format) (string, error) {
				return format.String(), nil
			},
		}),
		History: configReadHistory{
			Anilist: reg(field[bool, bool]{
				Key:         "read.history.anilist",
				Default:     true,
				Description: "Sync to Anilist reading history if logged in.",
			}),
			Local: reg(field[bool, bool]{
				Key:         "read.history.local",
				Default:     true,
				Description: "Save to local history.",
			}),
		},
		DownloadOnRead: reg(field[bool, bool]{
			Key:         "read.download_on_read",
			Default:     false,
			Description: "Download chapter to the default directory when opening for reading.",
		}),
	},
	Download: configDownload{
		Path: reg(field[string, string]{
			Key:         "download.path",
			Default:     xdg.UserDirs.Download,
			Description: "Path where chapters will be downloaded.",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
		Format: reg(field[string, libmangal.Format]{
			Key:         "download.format",
			Default:     libmangal.FormatPDF,
			Description: "Format to download chapters in.",
			Unmarshal: func(s string) (libmangal.Format, error) {
				return libmangal.FormatString(s)
			},
			Marshal: func(format libmangal.Format) (string, error) {
				return format.String(), nil
			},
		}),
		UserAgent: reg(field[string, string]{
			Key:         "download.user_agent",
			Default:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:126.0) Gecko/20100101 Firefox/126.0",
			Description: "User-Agent to use for making HTTP requests.",
		}),
		ModeDir: reg(field[int64, fs.FileMode]{
			Key:         "download.mode_dir_decimal",
			Default:     fs.FileMode(0o755),
			Description: "Permission bits used for all dirs created. Encodes into `int64` (decimal).",
			Unmarshal: func(i int64) (fs.FileMode, error) {
				return fs.FileMode(i), nil
			},
			Marshal: func(mode fs.FileMode) (int64, error) {
				return int64(mode), nil
			},
		}),
		ModeFile: reg(field[int64, fs.FileMode]{
			Key:         "download.mode_file_decimal",
			Default:     fs.FileMode(0o644),
			Description: "Permission bits used for all files created. Encodes into `int64` (decimal).",
			Unmarshal: func(i int64) (fs.FileMode, error) {
				return fs.FileMode(i), nil
			},
			Marshal: func(mode fs.FileMode) (int64, error) {
				return int64(mode), nil
			},
		}),
		ModeDB: reg(field[int64, fs.FileMode]{
			Key:         "download.mode_db_decimal",
			Default:     fs.FileMode(0o600),
			Description: "Permission bits used for database files created. Encodes into `int64` (decimal).",
			Unmarshal: func(i int64) (fs.FileMode, error) {
				return fs.FileMode(i), nil
			},
			Marshal: func(mode fs.FileMode) (int64, error) {
				return int64(mode), nil
			},
		}),
		Strict: reg(field[bool, bool]{
			Key:         "download.strict",
			Default:     true,
			Description: "If during metadata/banner/cover creation error occurs downloader will return it immediately and chapter won't be downloaded.",
		}),
		SkipIfExists: reg(field[bool, bool]{
			Key:         "download.skip_if_exists",
			Default:     true,
			Description: "Skip downloading chapter if its already downloaded (exists at path). Metadata will still be created if needed.",
		}),
		Provider: configDownloadProvider{
			CreateDir: reg(field[bool, bool]{
				Key:         "download.provider.create_dir",
				Default:     false,
				Description: "Create provider directory.",
			}),
			NameTemplate: reg(field[string, string]{
				Key:         "download.provider.name_template",
				Default:     "{{ .Name | sanitize }}",
				Description: "Template to use for naming downloaded providers.", // TODO: change these generic descriptions
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(funcs.FuncMap).
						Parse(s)
					return err
				},
			}),
		},
		Manga: configDownloadManga{
			CreateDir: reg(field[bool, bool]{
				Key:         "download.manga.create_dir",
				Default:     true,
				Description: "Create manga directory.",
			}),
			Cover: reg(field[bool, bool]{
				Key:         "download.manga.cover",
				Default:     false,
				Description: "Download manga cover.",
			}),
			Banner: reg(field[bool, bool]{
				Key:         "download.manga.banner",
				Default:     false,
				Description: "Download manga banner.",
			}),
			NameTemplate: reg(field[string, string]{
				Key:         "download.manga.name_template",
				Default:     `{{ printf "%s (%d) [alid-%d]" .AnilistManga.String .AnilistManga.StartDate.Year .AnilistManga.ID | sanitize }}`,
				Description: "Template to use for naming downloaded mangas, when Anilist data available.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(funcs.FuncMap).
						Parse(s)
					return err
				},
			}),
			NameTemplateFallback: reg(field[string, string]{
				Key:         "download.manga.name_template_fallback",
				Default:     `{{ .Title | sanitize }}`,
				Description: "Template to use for naming downloaded mangas, when no Anilist data is available.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(funcs.FuncMap).
						Parse(s)
					return err
				},
			}),
		},
		Volume: configDownloadVolume{
			CreateDir: reg(field[bool, bool]{
				Key:         "download.volume.create_dir",
				Default:     false,
				Description: "Create volume directory",
			}),
			NameTemplate: reg(field[string, string]{
				Key:         "download.volume.name_template",
				Default:     `{{ printf "Vol. %d" .Number | sanitize }}`,
				Description: "Template to use for naming downloaded volumes.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(funcs.FuncMap).
						Parse(s)
					return err
				},
			}),
		},
		Chapter: configDownloadChapter{
			NameTemplate: reg(field[string, string]{
				Key:         "download.chapter.name_template",
				Default:     `{{ printf "[%06.1f] %s" .Number .Title | sanitize }}`,
				Description: "Template to use for naming downloaded chapters.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(funcs.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Metadata: configDownloadMetadata{
			SeriesJSON: reg(field[bool, bool]{
				Key:         "download.metadata.series_json",
				Default:     true,
				Description: "Generate `series.json` file.",
			}),
			SkipSeriesJSONIfOngoing: reg(field[bool, bool]{
				Key:         "download.metadata.skip_series_json_if_ongoing",
				Default:     true,
				Description: "Avoid writing the `series.json` file if the manga is ongoing (publishing).",
			}),
			ComicInfoXML: reg(field[bool, bool]{
				Key:         "download.metadata.comicinfo_xml",
				Default:     true,
				Description: "Generate `ComicInfo.xml` file.",
			}),
		},
	},
	TUI: configTUI{
		ExpandSingleVolume: reg(field[bool, bool]{
			Key:         "tui.expand_single_volume",
			Default:     true,
			Description: "Skip selecting volume if there's only one.",
		}),
		Chapter: configTUIChapter{
			// TODO: add validation to the format
			NumberFormat: reg(field[string, string]{
				Key:         "tui.chapter.number_format",
				Default:     "[%06.1f]",
				Description: "Format that the chapter number (float 32) should take, for example '[%06.1f]'.",
			}),
			ShowNumber: reg(field[bool, bool]{
				Key:         "tui.chapter.show_number",
				Default:     true,
				Description: "If the chapter number should be shown prepended to the chapter title.",
			}),
			ShowDate: reg(field[bool, bool]{
				Key:         "tui.chapter.show_date",
				Default:     true,
				Description: "If the chapter date should be shown in the description.",
			}),
			ShowGroup: reg(field[bool, bool]{
				Key:         "tui.chapter.show_group",
				Default:     true,
				Description: "If the chapter scanlation group should be shown in the description.",
			}),
		},
	},
	Providers: configProviders{
		Path: reg(field[string, string]{
			Key:         "providers.path",
			Default:     filepath.Join(Dir, "providers"),
			Description: "Path where chapters will be placed/looked for.",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
		Parallelism: reg(field[int64, uint8]{
			Key:         "providers.parallelism",
			Default:     15,
			Description: "Parallelism to use for the scrapers that support it.",
		}),
		Headless: configProvidersHeadless{
			UseFlaresolverr: reg(field[bool, bool]{
				Key:         "providers.headless.use_flaresolverr",
				Default:     false,
				Description: "If Flaresolverr should be used for headless providers, requires a flaresolverr service to connect to. providers.headless.flaresolverr_url needs to be set.",
			}),
			FlaresolverrURL: reg(field[string, string]{
				Key:         "providers.headless.flaresolverr_url",
				Default:     "http://localhost:8191/v1",
				Description: "Flaresolverr URL to use for headless providers.",
			}),
		},
		Filter: configProvidersFilter{
			NSFW: reg(field[bool, bool]{
				Key:         "providers.filter.nsfw",
				Default:     true,
				Description: "If NSFW content should be included, usually used for MangaDex.",
			}),
			// TODO: add validation to language filter
			Language: reg(field[string, string]{
				Key:         "providers.filter.language",
				Default:     "en",
				Description: "The language the manga should be on.",
			}),
			MangaPlusQuality: reg(field[string, string]{
				Key:         "providers.filter.mangaplus_quality",
				Default:     "super_high",
				Description: `MangaPlus page image quality, one of "low", "high" or "super_high".`,
			}),
			MangaDexDataSaver: reg(field[bool, bool]{
				Key:         "providers.filter.mangadex_datasaver",
				Default:     false,
				Description: `Use MangaDex "data-saver" option for chapter pages.`,
			}),
			TitleChapterNumber: reg(field[bool, bool]{
				Key:         "providers.filter.title_chapter_number",
				Default:     false,
				Description: "Include the chapter number in the title regardless of the availability of the chapter title.",
			}),
			AvoidDuplicateChapters: reg(field[bool, bool]{
				Key:         "providers.filter.avoid_duplicate_chapters",
				Default:     true,
				Description: "Only select one chapter when multiple of the same number are present.",
			}),
			ShowUnavailableChapters: reg(field[bool, bool]{
				Key:         "providers.filter.show_unavailable_chapters",
				Default:     false,
				Description: "When there are non-downloadable chapters, show them anyways. Should only be used to search around.",
			}),
		},
		MangaPlus: configProvidersMangaPlus{
			OSVersion: reg(field[string, string]{
				Key:         "providers.mangaplus.os_version",
				Default:     "30",
				Description: "The OS Version used for the MangaPlus API calls.",
			}),
			AppVersion: reg(field[string, string]{
				Key:         "providers.mangaplus.app_version",
				Default:     "133",
				Description: "The App Version used for the MangaPlus API calls.",
			}),
			AndroidID: reg(field[string, string]{
				Key:         "providers.mangaplus.android_id",
				Default:     "",
				Description: "The Android ID used for the MangaPlus API calls. If empty will be randomly generated.",
			}),
		},
	},
	Library: configLibrary{
		Path: reg(field[string, string]{
			Key:         "library.path",
			Default:     "",
			Description: "Path to the manga library. Empty string will fallback to the download.path.",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
	},
}
