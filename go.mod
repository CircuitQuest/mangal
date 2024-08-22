module github.com/luevano/mangal

go 1.22

// Mistakenly "uploaded" when importing in a separate package...
retract (
	v0.0.0-20240707233308-9adbbbc770e1
	v0.0.0-20240517101323-1d20eb78430f
	v0.0.0-20240514061342-bf49b6004548
)

require (
	github.com/adrg/xdg v0.5.0
	github.com/charmbracelet/bubbles v0.19.0
	github.com/charmbracelet/bubbletea v0.27.0
	github.com/charmbracelet/lipgloss v0.13.0
	github.com/charmbracelet/x/ansi v0.2.3
	github.com/disgoorg/disgo v0.18.10
	github.com/fatih/camelcase v1.0.0
	github.com/getkin/kin-openapi v0.127.0
	github.com/go-git/go-git/v5 v5.12.0
	github.com/google/uuid v1.6.0
	github.com/itchyny/gojq v0.12.16
	github.com/json-iterator/go v1.1.12
	github.com/ktr0731/go-fuzzyfinder v0.8.0
	github.com/labstack/echo/v4 v4.12.0
	github.com/lithammer/fuzzysearch v1.1.8
	github.com/luevano/gopher-luadoc v0.3.2
	github.com/luevano/libmangal v0.20.1
	github.com/luevano/luaprovider v0.14.1
	github.com/luevano/mangoprovider v0.16.4
	github.com/muesli/reflow v0.3.0
	github.com/oapi-codegen/runtime v1.1.1
	github.com/pelletier/go-toml v1.9.5
	github.com/philippgille/gokv v0.7.0
	github.com/philippgille/gokv/bigcache v0.7.0
	github.com/philippgille/gokv/encoding v0.7.0
	github.com/philippgille/gokv/util v0.7.0
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.33.0
	github.com/samber/lo v1.47.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/afero v1.11.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/wk8/go-ordered-map/v2 v2.1.8
	github.com/yuin/gopher-lua v1.1.1
	github.com/zyedidia/generic v1.2.1
	go.etcd.io/bbolt v1.3.11
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa
	golang.org/x/oauth2 v0.22.0
)

require (
	github.com/Luzifer/go-openssl/v4 v4.2.2 // indirect
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/antchfx/htmlquery v1.3.2 // indirect
	github.com/antchfx/xmlquery v1.4.1 // indirect
	github.com/antchfx/xpath v1.3.1 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/charmbracelet/x/term v0.2.0 // indirect
	github.com/cyphar/filepath-securejoin v0.3.1 // indirect
	github.com/disgoorg/json v1.1.0 // indirect
	github.com/disgoorg/snowflake/v2 v2.0.3 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly/v2 v2.1.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/invopop/yaml v0.3.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/luevano/mangodex v0.3.8 // indirect
	github.com/luevano/mangoplus v0.4.4 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/saintfish/chardet v0.0.0-20230101081208-5e3ef4b5456d // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/temoto/robotstxt v1.1.2 // indirect
	github.com/tj/go-naturaldate v1.3.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/vineesh12344/gojsfuck v0.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

require (
	dario.cat/mergo v1.0.1 // indirect
	github.com/JohannesKaufmann/html-to-markdown v1.6.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProtonMail/go-crypto v1.0.0 // indirect
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/atotto/clipboard v0.1.4
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/charmbracelet/harmonica v0.2.0 // indirect
	github.com/cixtor/readability v1.0.0 // indirect
	github.com/cloudflare/circl v1.4.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/gdamore/tcell/v2 v2.7.4 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.5.0 // indirect
	github.com/go-rod/rod v0.116.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/hhrutter/lzw v1.0.0 // indirect
	github.com/hhrutter/tiff v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/ivanpirog/coloredcobra v1.0.1
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/ka-weihe/fast-levenshtein v0.0.0-20201227151214-4c99ee36a1ba // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/ktr0731/go-ansisgr v0.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.15.2
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/nsf/termbox-go v1.1.1 // indirect
	github.com/pdfcpu/pdfcpu v0.8.0 // indirect
	github.com/philippgille/gokv/syncmap v0.7.0 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/robertkrimen/otto v0.4.0 // indirect
	github.com/sahilm/fuzzy v0.1.1 // indirect
	github.com/segmentio/fasthash v1.0.3 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/ysmood/fetchup v0.2.4 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.40.0 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	github.com/ysmood/leakless v0.9.0 // indirect
	github.com/yuin/gluamapper v0.0.0-20150323120927-d836955830e7 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/image v0.19.0 // indirect
	golang.org/x/mod v0.20.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
