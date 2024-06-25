package anilistmangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(anilist *lmanilist.Anilist, mangas []lmanilist.Manga, onResponse onResponseFunc) *state {
	_keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"anilist manga", "anilist mangas",
		mangas,
		func(manga lmanilist.Manga) _list.DefaultItem {
			return &item{manga: manga}
		},
		_keyMap)

	sInput := textinput.New()
	sInput.Prompt = icon.Search.String() + " "
	// TODO: use the initial query as placeholder
	sInput.Placeholder = "Search anilist..."
	sInput.PromptStyle = style.Bold.Warning
	// TODO: better handle input limit
	sInput.CharLimit = 64

	s := &state{
		anilist:     anilist,
		list:        listWrapper,
		searchInput: sInput,
		searchState: unsearched,
		onResponse:  onResponse,
		keyMap:      _keyMap,
	}
	s.updateKeybindings()

	return s
}
