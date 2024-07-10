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
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list      *list.State
	loaded    *set.Set[*item]
	extraInfo *bool
	keyMap    keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return s.list.Intermediate()
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
	return base.Title{Text: "Providers"}
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
func (s *state) Update(ctx context.Context, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s.list.FilterState() == _list.Filtering {
			goto end
		}

		i, ok := s.list.SelectedItem().(*item)
		if !ok {
			return nil
		}

		switch {
		case key.Matches(msg, s.keyMap.confirm):
			return loadProviderCmd(i)
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
				item.markClosed()
			}

			return base.Notify("Closed all clients")
		}
	case loadProviderMsg:
		item := msg.item

		return tea.Sequence(
			base.Loading(fmt.Sprintf("Loading provider %q", item.loader)),
			func() tea.Msg {
				var mangalClient *libmangal.Client
				if c := client.Get(item.loader); c != nil {
					log.Log("Using existing mangal client for provider %q", item.loader)
					mangalClient = c
				} else {
					log.Log("New mangal client for provider %q", item.loader)
					c, err := client.NewClient(ctx, item.loader)
					if err != nil {
						return err
					}
					mangalClient = c
					item.markLoaded()

					mangalClient.Logger().SetOnLog(func(format string, a ...any) {
						log.Log(format, a...)
					})
				}
				return mangas.New(mangalClient)
			},
			base.Loaded,
		)
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
