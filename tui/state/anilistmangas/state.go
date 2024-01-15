package anilistmangas

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/listwrapper"
	"github.com/luevano/mangal/tui/state/loading"
	"github.com/luevano/mangal/tui/state/textinput"
)

var _ base.State = (*State)(nil)

type OnResponseFunc func(response *libmangal.AnilistManga) tea.Cmd

type State struct {
	anilist *libmangal.Anilist
	list    *listwrapper.State

	onResponse OnResponseFunc

	keyMap KeyMap
}

func (s *State) Intermediate() bool {
	return true
}

func (s *State) KeyMap() help.KeyMap {
	return s.keyMap
}

func (s *State) Title() base.Title {
	return base.Title{Text: "Anilist Mangas"}
}

func (s *State) Subtitle() string {
	return s.list.Subtitle()
}

func (s *State) Status() string {
	return s.list.Status()
}

func (s *State) Backable() bool {
	return s.list.Backable()
}

func (s *State) Resize(size base.Size) {
	s.list.Resize(size)
}

func (s *State) Update(model base.Model, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, s.keyMap.Confirm):
			item, ok := s.list.SelectedItem().(Item)
			if !ok {
				return nil
			}

			return s.onResponse(item.Manga)
		case key.Matches(msg, s.keyMap.Search):
			return func() tea.Msg {
				return textinput.New(textinput.Options{
					Title:        base.Title{Text: "Search Mangas"},
					Prompt:       "Enter Anilist manga title",
					Placeholder:  "",
					Intermediate: true,
					OnResponse: func(response string) tea.Cmd {
						return tea.Sequence(
							func() tea.Msg {
								return loading.New("Searching", fmt.Sprintf("Searching for %q on Anilist", response))
							},
							func() tea.Msg {
								mangas, err := s.anilist.SearchMangas(model.Context(), response)
								if err != nil {
									return err
								}

								return New(s.anilist, mangas, s.onResponse)
							},
						)
					},
				})
			}
		}
	}

end:
	return s.list.Update(model, msg)
}

func (s *State) View(model base.Model) string {
	return s.list.View(model)
}

func (s *State) Init(model base.Model) tea.Cmd {
	return s.list.Init(model)
}
