package mangas

import (
	_list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

func New(client *libmangal.Client) *state {
	_keyMap := newKeyMap()
	listWrapper := list.New(
		2,
		"manga", "mangas",
		nil,
		func(manga mangadata.Manga) _list.DefaultItem {
			return &item{manga}
		},
		_keyMap)

	sInput := textinput.New()
	sInput.Prompt = icon.Search.String() + " "
	// TODO: use the initial query as placeholder
	sInput.Placeholder = "Search manga..."
	sInput.PromptStyle = style.Bold.Warning
	// TODO: better handle input limit
	sInput.CharLimit = 64

	s := &state{
		list:        listWrapper,
		client:      client,
		searchInput: sInput,
		searchState: unsearched,
		keyMap:      _keyMap,
	}
	s.updateKeybindings()

	return s
}
