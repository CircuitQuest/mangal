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
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/mangas"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/state/wrapper/textinput"
	"github.com/pkg/errors"
)

var _ base.State = (*State)(nil)

// State implements base.State.
type State struct {
	list            *list.State
	providerLoaders []libmangal.ProviderLoader
	keyMap          keyMap
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
	return s.keyMap
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
func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

// Init implements base.State.
func (s *State) Init(ctx context.Context) tea.Cmd {
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

		item, ok := s.list.SelectedItem().(Item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return tea.Sequence(
				func() tea.Msg {
					return loading.New("Loading", fmt.Sprintf("Loading provider %q", item))
				},
				func() tea.Msg {
					mangalClient, err := client.NewClient(ctx, item)
					if err != nil {
						return err
					}

					mangalClient.Logger().SetOnLog(func(format string, a ...any) {
						log.Log(format, a...)
					})

					return textinput.New(textinput.Options{
						Title:       base.Title{Text: "Search Manga"},
						Subtitle:    fmt.Sprintf("Search using %q provider", mangalClient),
						Placeholder: "Manga title...",
						OnResponse: func(response string) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									return loading.New("Searching", fmt.Sprintf("Searching for %q", response))
								},
								func() tea.Msg {
									mangaList, err := mangalClient.SearchMangas(ctx, response)
									if err != nil {
										return err
									}

									return mangas.New(mangalClient, response, mangaList)
								},
							)
						},
					})
				},
			)
		case key.Matches(msg, s.keyMap.info):
			return func() tea.Msg {
				return errors.New("unimplemented")
			}
		}
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *State) View() string {
	return s.list.View()
}
