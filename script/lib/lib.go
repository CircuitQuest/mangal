package lib

import (
	luadoc "github.com/luevano/gopher-luadoc"
	"github.com/luevano/libmangal"
	sdk "github.com/luevano/luaprovider/lib"
	"github.com/luevano/mangal/afs"
	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/script/lib/client"
	"github.com/luevano/mangal/script/lib/json"
	"github.com/luevano/mangal/script/lib/prompt"
	lua "github.com/yuin/gopher-lua"
)

type Options struct {
	Client  *libmangal.Client
	Anilist *libmangal.Anilist
}

func Lib(state *lua.LState, options Options) *luadoc.Lib {
	SDKOptions := sdk.DefaultOptions()
	SDKOptions.FS = afs.Afero.Fs

	return &luadoc.Lib{
		Name:        meta.AppName,
		Description: meta.AppName + " scripting mode utilities",
		// TODO: addd anilist lib, why isn't it added?
		Libs: []*luadoc.Lib{
			sdk.Lib(state, SDKOptions),
			prompt.Lib(),
			json.Lib(),
			client.Lib(options.Client),
		},
	}
}

func Preload(state *lua.LState, options Options) {
	lib := Lib(state, options)
	state.PreloadModule(lib.Name, lib.Loader())
}
