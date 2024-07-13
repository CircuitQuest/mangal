package metadata

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/color"
	_style "github.com/luevano/mangal/theme/style"
)

// TODO: handle wrapping for smaller terminals
func (s *State) renderMetadata() string {
	m := s.meta.Metadata()
	f := allFields(m)

	s.enumeratorStyle = lipgloss.NewStyle().
		Foreground(s.meta.Style().Color).
		MarginRight(1).
		Padding(0, 0, 0, 2)

	// prone to errros
	l := make([]string, 12)
	l[0] = jH(s.render(f.id), s.render(f.title))
	l[1] = jH(s.render(f.status), s.render(f.chapters), s.render(f.startDate), s.render(f.endDate))
	l[2] = jH(s.render(f.country), s.render(f.score), s.render(f.format), s.render(f.publisher))
	l[3] = s.render(f.url)
	l[4] = s.render(f.description)
	l[5] = s.render(f.cover)
	l[6] = s.render(f.banner)
	l[7] = s.render(f.characters)
	l[8] = jH(jV(s.render(f.genres), s.render(f.tags)), s.render(f.alternateTitles))
	l[9] = jH(s.render(f.authors), s.render(f.artists), s.render(f.translators), s.render(f.letterers))
	l[10] = s.render(f.extraIDs)
	l[11] = s.render(f.notes)

	return jV(l...)
}

func (s *State) render(f field) string {
	var str string
	switch value := f.value.(type) {
	case string:
		str = value
	case []string:
		return s.renderList(f.name, value, lipgloss.NewStyle().Width(f.width).Foreground(color.Secondary))
	case int:
		if value != 0 {
			str = strconv.Itoa(value)
		} else {
			str = ""
		}
	case float64:
		str = strconv.FormatFloat(value, 'f', -1, 32)
	case float32:
		str = strconv.FormatFloat(float64(value), 'f', -1, 32)
	case metadata.Status:
		str = string(value)
	case metadata.Date:
		str = value.String()
	case metadata.ID:
		if value.Code != "" {
			name := f.name + "(" + value.Code + ")"
			return s.renderField(name, str, f.width)
		}
		str = value.Raw
	case []metadata.ID:
		// convert IDS to rendered strings
		strs := make([]string, len(value))
		for i, id := range value {
			style := s.meta.IDStyle(id.Source)
			strs[i] = lipgloss.NewStyle().Foreground(style.Color).Render(style.Prefix+": ") + value[i].Raw
		}
		return s.renderList(f.name, strs, lipgloss.NewStyle().Width(f.width))
	default:
		str = _style.Bold.Error.Render(fmt.Sprintf("RenderField: unkown field type: %T", value))
	}
	return s.renderField(f.name, str, f.width)
}

func (s *State) renderFieldName(name string) string {
	// TODO: use icon.Field once the reset style is fixed where
	// once a style is applied, a style on top will be cut off
	// by the first one (due to the "reset" ansi code)
	return lipgloss.NewStyle().
		Foreground(s.meta.Style().Color).
		Render(">" + name + ":")
}

func (s *State) renderField(name string, str string, width int) string {
	minWidth := 2
	if width < minWidth {
		width = minWidth
	}
	if str == "" {
		str = _style.Normal.Secondary.Render(strings.Repeat(".", width))
	}
	return jH(s.renderFieldName(name), lipgloss.NewStyle().Width(width).MaxWidth(width).UnsetMaxWidth().Render(str))
}

func (s *State) renderList(name string, items []string, itemStyle lipgloss.Style) string {
	if len(items) == 0 {
		return s.renderField(name, "", 3)
	}
	l := list.New(items).
		Enumerator(list.Dash).
		EnumeratorStyle(s.enumeratorStyle).
		ItemStyle(itemStyle)
	// jV would insert a space inbetween
	return lipgloss.JoinVertical(lipgloss.Left, s.renderFieldName(name), l.String())
}

func jH(strs ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, insertSpacesInbetween(strs)...)
}

func jV(strs ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, insertSpacesInbetween(strs)...)
}

func insertSpacesInbetween(strs []string) []string {
	size := len(strs)
	newStrs := make([]string, 2*size-1)
	for i := range strs {
		newStrs[2*i] = strs[i]
		if i < size-1 {
			newStrs[2*i+1] = " "
		}
	}
	return newStrs
}
