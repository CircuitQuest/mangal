package metadata

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/tui/util"
)

var _ tea.Model = (*Model)(nil)

// Model implements tea.Model.
type Model struct {
	meta  metadata.Metadata
	style lipgloss.Style

	metaStyle util.MetaStyle
	ShowFull  bool
}

// SetMetadata replaces the current metadata and updates the style.
func (m *Model) SetMetadata(meta metadata.Metadata) {
	m.meta = meta
	m.updateStyle()
}

// updateStyle sets the style based on the current metadata.
func (m *Model) updateStyle() {
	m.metaStyle = util.MetaIDStyle(m.meta.ID())
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View implements tea.Model.
func (m *Model) View() string {
	var y string
	if m.meta.StartDate().Year != 0 {
		y = " (" + strconv.Itoa(m.meta.StartDate().Year) + ")"
	}

	if !m.ShowFull {
		return m.style.Background(m.metaStyle.Color).Render(m.meta.Title() + y)
	}

	i := ": "
	if id := m.meta.ID().Raw; id != "" {
		i = " (" + id + "): "
	}

	p := m.meta.ID().Code
	if p != "" {
		p = "[" + p + "] "
	}
	return m.style.Background(m.metaStyle.Color).Render(p + m.metaStyle.Prefix + i + m.meta.Title() + y)
}
