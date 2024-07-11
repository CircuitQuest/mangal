package chapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	_list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/anilist"
	"github.com/luevano/mangal/tui/state/formats"
	"github.com/luevano/mangal/tui/state/wrapper/list"
	"github.com/zyedidia/generic/set"
)

var _ base.State = (*state)(nil)

// state implements base.state.
type state struct {
	list *list.State

	chapters []mangadata.Chapter
	volume   mangadata.Volume // can be nil
	manga    mangadata.Manga
	client   *libmangal.Client

	selected set.Set[*item]

	renderedSep             string
	renderedSubtitleFormats string

	// to avoid extra read/download cmds from
	// firing up when an action is already happening,
	// only blocks keymaps handled by Update for read/download
	actionRunning string

	showVolumeNumber  *bool
	showChapterNumber *bool
	showGroup         *bool
	showDate          *bool

	styles styles
	keyMap keyMap
}

// Intermediate implements base.State.
func (s *state) Intermediate() bool {
	return false
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
	return base.Title{Text: s.manga.String()}
}

// Subtitle implements base.State.
func (s *state) Subtitle() string {
	var subtitle strings.Builder
	subtitle.Grow(100)

	subtitle.WriteString(s.list.Subtitle())
	if s.selected.Size() > 0 {
		selected := s.renderedSep +
			s.styles.subtitle.Render(fmt.Sprintf("%d selected", s.selected.Size()))
		subtitle.WriteString(selected)
	}
	subtitle.WriteString(s.renderedSubtitleFormats)

	return subtitle.String()
}

// Status implements base.State.
func (s *state) Status() string {
	if s.volume != nil {
		return fmt.Sprintf("Vol. %s%s%s", s.volume, s.renderedSep, s.list.Status())
	}
	return s.list.Status()
}

// Resize implements base.State.
func (s *state) Resize(size base.Size) tea.Cmd {
	return s.list.Resize(size)
}

// Init implements base.State.
func (s *state) Init(ctx context.Context) tea.Cmd {
	s.updateRenderedSubtitleFormats()
	return s.list.Init(ctx)
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
		case key.Matches(msg, s.keyMap.toggle):
			i.toggle()
			if i.selected {
				s.selected.Put(i)
			} else {
				s.selected.Remove(i)
			}
			return nil
		case key.Matches(msg, s.keyMap.unselectAll):
			for _, item := range s.selected.Keys() {
				item.toggle()
				s.selected.Remove(item)
			}
			return nil
		case key.Matches(msg, s.keyMap.selectAll):
			for _, listItem := range s.list.Items() {
				it, ok := listItem.(*item)
				if !ok {
					continue
				}

				if !it.selected {
					it.toggle()
					s.selected.Put(it)
				}
			}
			return nil
		case key.Matches(msg, s.keyMap.changeFormat):
			return func() tea.Msg {
				return formats.New()
			}
		case key.Matches(msg, s.keyMap.openURL):
			return s.openURLCmd(i.chapter)
		case key.Matches(msg, s.keyMap.download):
			return s.downloadCmd(ctx, i)
		case key.Matches(msg, s.keyMap.read):
			return s.readCmd(ctx, i)
		case key.Matches(msg, s.keyMap.anilist):
			return func() tea.Msg {
				return anilist.New(s.client.Anilist(), s.manga)
			}
		case key.Matches(msg, s.keyMap.toggleVolumeNumber):
			*s.showVolumeNumber = !(*s.showVolumeNumber)
		case key.Matches(msg, s.keyMap.toggleChapterNumber):
			*s.showChapterNumber = !(*s.showChapterNumber)
		case key.Matches(msg, s.keyMap.toggleGroup):
			*s.showGroup = !(*s.showGroup)
			s.updateListDelegate()
		case key.Matches(msg, s.keyMap.toggleDate):
			*s.showDate = !(*s.showDate)
			s.updateListDelegate()
		}
	case base.RestoredMsg:
		// usually the downloaded chapters change or the metadata when restoring the chapter list
		s.updateAllItems()
		s.updateRenderedSubtitleFormats()
	}
end:
	return s.list.Update(ctx, msg)
}

// View implements base.State.
func (s *state) View() string {
	return s.list.View()
}
