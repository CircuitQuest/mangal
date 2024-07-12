package metadata

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/libmangal/metadata"
)

var _ tea.Model = (*Model)(nil)

// Model implements tea.Model.
type Model struct {
	meta metadata.Metadata

	styles styles
	style  lipgloss.Style
	prefix string

	ShowFull bool
}

// SetMetadata replaces the current metadata and updates the style.
func (m *Model) SetMetadata(meta metadata.Metadata) {
	m.meta = meta
	m.updateStyle()
}

// updateStyle sets the style based on the current metadata.
func (m *Model) updateStyle() {
	switch m.meta.ID().IDSource {
	case metadata.IDSourceAnilist:
		m.style = m.styles.anilist
		m.prefix = "Anilist"
	default:
		m.style = m.styles.provider
		m.prefix = "Provider"
	}
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
	year := " (" + strconv.Itoa(m.meta.StartDate().Year) + ")"
	if !m.ShowFull {
		return m.style.Render(m.meta.Title() + year)
	}

	sep := ": "
	if id := m.meta.ID().IDRaw; id != "" {
		sep = " (" + id + "): "
	}
	return m.style.Render(m.prefix + sep + m.meta.Title() + year)
}
