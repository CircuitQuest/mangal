package config

import (
	"text/template"
	"time"

	"github.com/adrg/xdg"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/icon"
	"github.com/luevano/mangal/template/util"
)

var Config = config{
	Icons: reg(Field[string, icon.Type]{
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
	CLI: configCLI{
		ColoredHelp: reg(Field[bool, bool]{
			Key:         "cli.colored_help",
			Default:     true,
			Description: "Enable colors in cli help.",
		}),
		Mode: configCLIMode{
			// TODO: change mode to either none or some (new) "inline" mode
			Default: reg(Field[string, Mode]{
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
		Format: reg(Field[string, libmangal.Format]{
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
			Anilist: reg(Field[bool, bool]{
				Key:         "read.history.anilist",
				Default:     true,
				Description: "Sync to Anilist reading history if logged in.",
			}),
			Local: reg(Field[bool, bool]{
				Key:         "read.history.local",
				Default:     true,
				Description: "Save to local history.",
			}),
		},
		DownloadOnRead: reg(Field[bool, bool]{
			Key:         "read.download_on_read",
			Default:     false,
			Description: "Download chapter to the default directory when opening for reading.",
		}),
	},
	Download: configDownload{
		// Don't use config.Config.Download.Path.Get()
		// as it creates a directory when called, may be unwanted?
		Path: reg(Field[string, string]{
			Key:         "download.path",
			Default:     xdg.UserDirs.Download,
			Description: "Path where chapters will be downloaded.",
			Unmarshal: func(s string) (string, error) {
				return expandPath(s)
			},
		}),
		Format: reg(Field[string, libmangal.Format]{
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
		Strict: reg(Field[bool, bool]{
			Key:         "download.strict",
			Default:     true,
			Description: "If during metadata/banner/cover creation error occurs downloader will return it immediately and chapter won't be downloaded.",
		}),
		SkipIfExists: reg(Field[bool, bool]{
			Key:         "download.skip_if_exists",
			Default:     true,
			Description: "Skip downloading chapter if its already downloaded (exists at path). Metadata will still be created if needed.",
		}),
		Manga: configDownloadManga{
			CreateDir: reg(Field[bool, bool]{
				Key:         "download.manga.create_dir",
				Default:     true,
				Description: "Create manga directory.",
			}),
			Cover: reg(Field[bool, bool]{
				Key:         "download.manga.cover",
				Default:     false,
				Description: "Download manga cover.",
			}),
			Banner: reg(Field[bool, bool]{
				Key:         "download.manga.banner",
				Default:     false,
				Description: "Download manga banner.",
			}),
			// TODO: in the future change this to a standardized name, maybe coming from anilist
			NameTemplate: reg(Field[string, string]{
				Key:         "download.manga.name_template",
				Default:     `{{ .Title | sanitize }}`,
				Description: "Template to use for naming downloaded mangas.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Volume: configDownloadVolume{
			CreateDir: reg(Field[bool, bool]{
				Key:         "download.volume.create_dir",
				Default:     false,
				Description: "Create volume directory",
			}),
			NameTemplate: reg(Field[string, string]{
				Key:         "download.volume.name_template",
				Default:     `{{ printf "Vol. %d" .Number | sanitize }}`,
				Description: "Template to use for naming downloaded volumes.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Chapter: configDownloadChapter{
			NameTemplate: reg(Field[string, string]{
				Key:         "download.chapter.name_template",
				Default:     `{{ printf "[%06.1f] %s" .Number .Title | sanitize }}`,
				Description: "Template to use for naming downloaded chapters.",
				Validate: func(s string) error {
					_, err := template.
						New("").
						Funcs(util.FuncMap).
						Parse(s)

					return err
				},
			}),
		},
		Metadata: configDownloadMetadata{
			ComicInfoXML: reg(Field[bool, bool]{
				Key:         "download.metadata.comicinfo_xml",
				Default:     true,
				Description: "Generate `ComicInfo.xml` file.",
			}),
			SeriesJSON: reg(Field[bool, bool]{
				Key:         "download.metadata.series_json",
				Default:     true,
				Description: "Generate `series.json` file.",
			}),
		},
	},
	TUI: configTUI{
		ExpandSingleVolume: reg(Field[bool, bool]{
			Key:         "tui.expand_single_volume",
			Default:     true,
			Description: "Skip selecting volume if there's only one.",
		}),
		Chapter: configTUIChapter{
			// TODO: add validation to the format
			NumberFormat: reg(Field[string, string]{
				Key:         "tui.chapter.number_format",
				Default:     "[%06.1f]",
				Description: "Format that the chapter number (float 32) should take, for example '[%06.1f]'.",
			}),
			ShowNumber: reg(Field[bool, bool]{
				Key:         "tui.chapter.show_number",
				Default:     true,
				Description: "If the chapter number should be shown prepended to the chapter title.",
			}),
			ShowDate: reg(Field[bool, bool]{
				Key:         "tui.chapter.show_date",
				Default:     true,
				Description: "If the chapter date should be shown in the description.",
			}),
			ShowGroup: reg(Field[bool, bool]{
				Key:         "tui.chapter.show_group",
				Default:     true,
				Description: "If the chapter scanlation group should be shown in the description.",
			}),
		},
	},
	Providers: configProviders{
		Parallelism: reg(Field[uint8, uint8]{
			Key:         "providers.parallelism",
			Default:     15,
			Description: "Parallelism to use for the scrapers that support it.",
		}),
		Cache: configProvidersCache{
			TTL: reg(Field[string, string]{
				Key:         "providers.cache.ttl",
				Default:     "24h",
				Description: `Time to live. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".`,
				Validate: func(s string) error {
					_, err := time.ParseDuration(s)
					return err
				},
			}),
		},
		Headless: configProvidersHeadless{
			UseFlaresolverr: reg(Field[bool, bool]{
				Key:         "providers.headless.use_flaresolverr",
				Default:     false,
				Description: "If Flaresolverr should be used for headless providers, requires a flaresolverr service to connect to. providers.headless.flaresolverr_url needs to be set.",
			}),
			FlaresolverrURL: reg(Field[string, string]{
				Key:         "providers.headless.flaresolverr_url",
				Default:     "http://localhost:8191/v1",
				Description: "Flaresolverr URL to use for headless providers.",
			}),
		},
	},
	Library: configLibrary{
		Path: reg(Field[string, string]{
			Key:         "library.path",
			Default:     "",
			Description: "Path to the manga library. Empty string will fallback to the download.path.",
		}),
	},
	Filter: configFilter{
		NSFW: reg(Field[bool, bool]{
			Key:         "filter.nsfw",
			Default:     true,
			Description: "If NSFW content should be included, usually used for MangaDex.",
		}),
		// TODO: add validation to language filter
		Language: reg(Field[string, string]{
			Key:         "filter.language",
			Default:     "en",
			Description: "The language the manga should be on.",
		}),
		MangaDexDataSaver: reg(Field[bool, bool]{
			Key:         "filter.mangadex_datasaver",
			Default:     false,
			Description: `Use MangaDex "data-saver" option for chapter pages.`,
		}),
		TitleChapterNumber: reg(Field[bool, bool]{
			Key:         "filter.title_chapter_number",
			Default:     false,
			Description: "Include the chapter number in the title regardless of the availability of the chapter title.",
		}),
		AvoidDuplicateChapters: reg(Field[bool, bool]{
			Key:         "filter.avoid_duplicate_chapters",
			Default:     true,
			Description: "Only select one chapter when multiple of the same number are present.",
		}),
		ShowUnavailableChapters: reg(Field[bool, bool]{
			Key:         "filter.show_unavailable_chapters",
			Default:     false,
			Description: "When there are non-downloadable chapters, show them anyways. Should only be used to search around.",
		}),
	},
}
