package providers

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/client"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/mangas"
	"github.com/luevano/mangal/tui/state/textinput"
	"github.com/mangalorg/libmangal"
	"github.com/pkg/errors"
)

var _ base.State = (*State)(nil)

type State struct {
	providersLoaders []libmangal.ProviderLoader
	list             *listwrapper.State
	keyMap           KeyMap
}

// Backable implements base.State.
func (s *State) Backable() bool {
	return s.list.Backable()
}

// Init implements base.State.
func (s *State) Init(model base.Model) tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return client.CloseAll()
		},
		s.list.Init(model),
	)
}

// Intermediate implements base.State.
func (s *State) Intermediate() bool {
	return s.list.Intermediate()
}

// KeyMap implements base.State.
func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

// Resize implements base.State.
func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

// Status implements base.State.
func (s *State) Status() string {
	return s.list.Status()
}

// Title implements base.State.
func (s *State) Title() base.Title {
	return base.Title{Text: "Providers"}
}

func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

// Update implements base.State.
func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
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
					return loading.New("Loading...", "")
				},
				func() tea.Msg {
					client, err := client.NewClient(model.Context(), item)
					if err != nil {
						return err
					}

					return textinput.New(textinput.Options{
						Title:  base.Title{Text: "Search"},
						Prompt: fmt.Sprintf("Using %q provider", client),
						OnResponse: func(response string) tea.Cmd {
							return tea.Sequence(
								func() tea.Msg {
									return loading.New("Loading", fmt.Sprintf("Searching for %q", response))
								},
								func() tea.Msg {
									m, err := client.SearchMangas(model.Context(), response)
									if err != nil {
										return err
									}

									return mangas.New(client, response, m)
								},
							)
						},
					})
				},
			)
		case key.Matches(msg, s.keyMap.info):
			return func() tea.Msg {
				return errors.New("not implemented")
			}
		}
	}
end:
	return s.list.Update(model, msg)
}

// View implements base.State.
func (s *State) View(model base.Model) string {
	return s.list.View(model)
}
