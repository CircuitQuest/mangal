package anilistmangas

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/model/search"
	"github.com/luevano/mangal/tui/state/wrapper/list"
)

type onResponseFunc func(manga lmanilist.Manga) tea.Cmd

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list    *list.State
	search  *search.Model
	anilist *lmanilist.Anilist

	onResponse onResponseFunc

	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
}

// Backable implements base.State.
func (s *state) Backable() bool {
	return s.list.Backable() && s.search.State() != search.Searching
}

// KeyMap implements base.State.
func (s *state) KeyMap() help.KeyMap {
	return s.list.KeyMap()
}

// Title implements base.State.
func (s *state) Title() base.Title {
	return base.Title{
		Text:       "Anilist Mangas",
		Background: color.Anilist.Background,
		Foreground: color.Anilist.Foreground,
	}
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
	s.search.Resize(size)
	sSize := base.Size{}
	sSize.Width, sSize.Height = lipgloss.Size(s.search.View())

	final := size
	final.Height -= sSize.Height

	return s.list.Resize(final)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	return s.list.Init(ctx)
}

// Updateimplements base.State.
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case search.SearchMsg:
		query := string(msg)

		return tea.Sequence(
			base.Loading(fmt.Sprintf("Searching %q on Anilist", query)),
			func() tea.Msg {
				mangas, err := s.anilist.SearchMangas(ctx, query)
				if err != nil {
					return err
				}

				items := make([]_list.Item, len(mangas))

				for i, m := range mangas {
					items[i] = &item{manga: m}
				}

				s.list.SetItems(items)
				return nil
			},
			base.Loaded,
		)
	}

	if s.search.State() == search.Searching {
		input, updateCmd := s.search.Update(msg)
		s.search = input.(*search.Model)
		return updateCmd
	}
	return s.handleBrowsingCmd(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.search.View(),
		s.list.View(),
	)
}
