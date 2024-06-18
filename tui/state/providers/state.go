package providers

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/mangas"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/state/wrapper/textinput"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	list      *list.State
	loaded    *set.Set[*Item]
	extraInfo *bool
	keyMap    keyMap
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return s.list.Intermediate()
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.list.Backable()
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "Providers"}
}

// Subtitle implements base.State.
func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *State) Status() string {
	return s.list.Status()
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
	// TODO: decide if Init should close all clients, instead use
	// State.Destroy() (if implemented) method and perform that there?
	return tea.Sequence(
		func() tea.Msg {
			return client.CloseAll()
		},
		s.list.Init(ctx),
	)
}

// Update implements base.State.
func (s *State) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		item, ok := s.list.SelectedItem().(*Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			var mangalClient *libmangal.Client

			return tea.Sequence(
				base.Loading(fmt.Sprintf("Loading provider %q", item)),
				func() tea.Msg {
					if c := client.Get(item.ProviderLoader); c != nil {
						log.Log("Using existing mangal client for provider %q", item)
						mangalClient = c
					} else {
						log.Log("New mangal client for provider %q", item)
						c, err := client.NewClient(ctx, item.ProviderLoader)
						if err != nil {
							return err
						}
						mangalClient = c
						item.MarkLoaded()

						mangalClient.Logger().SetOnLog(func(format string, a ...any) {
							log.Log(format, a...)
						})
					}
					return nil
				},
				base.Loaded,
				func() tea.Msg {
					return textinput.New(textinput.Options{
						Title:       base.Title{Text: "Search Manga"},
						Subtitle:    fmt.Sprintf("Search using %q provider", mangalClient),
						Placeholder: "Manga title...",
						OnResponse: func(response string) tea.Cmd {
							return tea.Sequence(
								base.Loading(fmt.Sprintf("Searching for %q", response)),
								func() tea.Msg {
									mangaList, err := mangalClient.SearchMangas(ctx, response)
									if err != nil {
										return err
									}

									return mangas.New(mangalClient, response, mangaList)
								},
								base.Loaded,
							)
						},
					})
				},
			)
		case key.Matches(msg, s.keyMap.info):
			*s.extraInfo = !(*s.extraInfo)

			if *s.extraInfo {
				s.list.SetDelegateHeight(3)
			} else {
				s.list.SetDelegateHeight(2)
			}
		case key.Matches(msg, s.keyMap.closeAll):
			if err := client.CloseAll(); err != nil {
				return func() tea.Msg {
					return err
				}
			}

			for _, item := range s.loaded.Keys() {
				item.MarkClosed()
			}

			return base.Notify("Closed all clients")
		}
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *State) View() string {
	return s.list.View()
}
