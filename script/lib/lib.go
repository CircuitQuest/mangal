package lib

import (
	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/script/lib/client"
	"github.com/luevano/mangal/script/lib/json"
	"github.com/luevano/mangal/script/lib/prompt"
	luadoc "github.com/mangalorg/gopher-luadoc"
	"github.com/mangalorg/libmangal"
	luaprovidersdk "github.com/mangalorg/luaprovider/lib"
	lua "github.com/yuin/gopher-lua"
)

const libName = meta.AppName

type Options struct {
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}

func Lib(state *lua.LState, options Options) *luadoc.Lib {
	SDKOptions := luaprovidersdk.DefaultOptions()
	SDKOptions.FS = afs.Afero.Fs

	lib := &luadoc.Lib{
		Name:        libName,
		Description: meta.AppName + " scripting mode utilities",
		Libs: []*luadoc.Lib{
			luaprovidersdk.Lib(state, SDKOptions),
			prompt.Lib(),
			json.Lib(),
			client.Lib(options.Client),
		},
	}

	return lib
}

func Preload(state *lua.LState, options Options) {
	lib := Lib(state, options)
	state.PreloadModule(lib.Name, lib.Loader())
}
