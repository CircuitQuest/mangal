package anilist

import (
	"context"

	luadoc "github.com/luevano/gopher-luadoc"
	"github.com/luevano/libmangal/metadata"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/script/lib/util"
	lua "github.com/yuin/gopher-lua"
)

const (
	libName = "anilist"

	mangaTypeName = libName + "_manga"
)

func Lib(anilist *metadata.ProviderWithCache) *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "Anilist operations",
		Funcs: []*luadoc.Func{
			{
				Name:        "search_mangas",
				Description: "Search mangas on Anilist",
				Value:       newSearchMangas(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "query",
						Description: "Query to search",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "mangas",
						Description: "Anilist mangas",
						Type:        luadoc.List(mangaTypeName),
					},
				},
			},
			{
				Name:        "find_closest_mangas",
				Description: "",
				Value:       newFindClosestManga(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "title",
						Description: "Manga title to search for",
						Type:        luadoc.String,
					},
				},
				Returns: []*luadoc.Param{
					{
						Name:        "manga",
						Description: "Closest manga found",
						Type:        mangaTypeName,
					},
					{
						Name:        "found",
						Description: "Whether the closest manga was found",
						Type:        luadoc.Boolean,
					},
				},
			},
			{
				Name:        "bind_title_with_id",
				Description: "Binds manga title to anilist manga id",
				Value:       newBindTitleWithID(anilist),
				Params: []*luadoc.Param{
					{
						Name:        "title",
						Description: "Manga title to use for binding",
						Type:        luadoc.String,
					},
					{
						Name:        "id",
						Description: "Anilist manga ID to bind to",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}

func notAnilistProvider(state *lua.LState) int {
	state.RaiseError("anilist client is missing")
	return 0
}

func newSearchMangas(anilist *metadata.ProviderWithCache) lua.LGFunction {
	if anilist.Info().Source != metadata.IDSourceAnilist {
		return notAnilistProvider
	}

	return func(state *lua.LState) int {
		query := state.CheckString(1)

		metas, err := anilist.Search(state.Context(), query)
		util.Must(state, err)

		// TODO: better handle this?
		mangas := make([]lmanilist.Manga, len(metas))
		for i, m := range metas {
			manga := m.(*lmanilist.Manga)
			mangas[i] = *manga
		}

		table := util.SliceToTable(state, mangas, func(manga lmanilist.Manga) lua.LValue {
			return util.NewUserData(state, manga, mangaTypeName)
		})

		state.Push(table)
		return 1
	}
}

func newFindClosestManga(anilist *metadata.ProviderWithCache) lua.LGFunction {
	if anilist.Info().Source != metadata.IDSourceAnilist {
		return notAnilistProvider
	}

	return func(state *lua.LState) int {
		title := state.CheckString(1)

		manga, found, err := anilist.FindClosest(context.Background(), title, 3, 3)
		util.Must(state, err)

		util.Push(state, manga, mangaTypeName)
		state.Push(lua.LBool(found))
		return 2
	}
}

func newBindTitleWithID(anilist *metadata.ProviderWithCache) lua.LGFunction {
	if anilist.Info().Source != metadata.IDSourceAnilist {
		return notAnilistProvider
	}

	return func(state *lua.LState) int {
		title := state.CheckString(1)
		ID := state.CheckInt(2)

		err := anilist.BindTitleWithID(title, ID)
		util.Must(state, err)

		return 0
	}
}
