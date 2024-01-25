package script

import (
	"context"
	"io"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/client/anilist"
	"github.com/luevano/mangal/provider/loader"
	"github.com/luevano/mangal/script/lib"
	lua "github.com/yuin/gopher-lua"
)

type Variables = map[string]string

type Args struct {
	File          string
	String        string
	Stdin         bool
	Provider      string
	Variables     Variables
	LoaderOptions *loader.Options
}

type Options struct {
	Client    *libmangal.Client
	Anilist   *libmangal.Anilist
	Variables Variables
}

func addVarsTable(state *lua.LState, variables Variables) {
	table := state.NewTable()
	for key, value := range variables {
		table.RawSetString(key, lua.LString(value))
	}

	state.SetGlobal("Vars", table)
}

func addLibraries(state *lua.LState, options lib.Options) {
	lib.Preload(state, options)
}

func Run(ctx context.Context, args Args, script io.Reader) error {
	client, err := client.NewClientByID(context.Background(), args.Provider, *args.LoaderOptions)
	if err != nil {
		return err
	}

	state := lua.NewState()
	state.SetContext(ctx)

	addVarsTable(state, args.Variables)
	addLibraries(state, lib.Options{
		Client:  client,
		Anilist: anilist.Anilist,
	})

	lFunction, err := state.Load(script, "script")
	if err != nil {
		return err
	}

	return state.CallByParam(lua.P{
		Fn:      lFunction,
		NRet:    1,
		Protect: true,
	})
}
