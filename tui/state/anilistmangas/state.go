package anilistmangas

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/luevano/mangal/tui/state/wrapper/textinput"
)

var _ base.State = (*state)(nil)

type onResponseFunc func(response *lmanilist.Manga) tea.Cmd

// state implements base.state.
type state struct {
	anilist *lmanilist.Anilist
	list    *list.State

	onResponse onResponseFunc

	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return true
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.Backable()
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{Text: "Anilist Mangas"}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	return s.list.Subtitle()
}

// Status implements base.State.
func (s *state) Status() string {
	return s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	return s.list.Init(ctx)
}

// Updateimplements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			i, ok := s.list.SelectedItem().(*item)
			if !ok {
				return nil
			}

			return s.onResponse(i.manga)
		case key.Matches(msg, s.keyMap.search):
			return func() tea.Msg {
				return textinput.New(textinput.Options{
					Title:        base.Title{Text: "Search Anilist"},
					Subtitle:     "Search Anilist manga",
					Placeholder:  "Anilist manga title...",
					Intermediate: true,
					OnResponse: func(response string) tea.Cmd {
						return tea.Sequence(
							base.Loading(fmt.Sprintf("Searching %q on Anilist", response)),
							func() tea.Msg {
								mangas, err := s.anilist.SearchMangas(ctx, response)
								if err != nil {
									return err
								}

								return New(s.anilist, mangas, s.onResponse)
							},
							base.Loaded,
						)
					},
				})
			}
		}
	}

end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
