package libs

import (
	"github.com/luevano/mangal-lua-libs/base64"
	"github.com/luevano/mangal-lua-libs/crypto"
	"github.com/luevano/mangal-lua-libs/filepath"
	"github.com/luevano/mangal-lua-libs/goos"
	"github.com/luevano/mangal-lua-libs/headless"
	"github.com/luevano/mangal-lua-libs/html"
	"github.com/luevano/mangal-lua-libs/http"
	"github.com/luevano/mangal-lua-libs/humanize"
	"github.com/luevano/mangal-lua-libs/inspect"
	"github.com/luevano/mangal-lua-libs/ioutil"
	"github.com/luevano/mangal-lua-libs/json"
	"github.com/luevano/mangal-lua-libs/log"
	"github.com/luevano/mangal-lua-libs/regexp"
	"github.com/luevano/mangal-lua-libs/runtime"
	"github.com/luevano/mangal-lua-libs/shellescape"
	"github.com/luevano/mangal-lua-libs/stats"
	"github.com/luevano/mangal-lua-libs/storage"
	"github.com/luevano/mangal-lua-libs/strings"
	"github.com/luevano/mangal-lua-libs/template"
	"github.com/luevano/mangal-lua-libs/time"
	"github.com/luevano/mangal-lua-libs/xmlpath"
	"github.com/luevano/mangal-lua-libs/yaml"
	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher lua packages
func Preload(L *lua.LState) {
	for _, preload := range []func(*lua.LState){
		yaml.Preload,
		html.Preload,
		headless.Preload,
		xmlpath.Preload,
		time.Preload,
		template.Preload,
		strings.Preload,
		storage.Preload,
		stats.Preload,
		shellescape.Preload,
		runtime.Preload,
		regexp.Preload,
		log.Preload,
		json.Preload,
		ioutil.Preload,
		inspect.Preload,
		humanize.Preload,
		http.Preload,
		goos.Preload,
		filepath.Preload,
		crypto.Preload,
		base64.Preload,
	} {
		preload(L)
	}
}
