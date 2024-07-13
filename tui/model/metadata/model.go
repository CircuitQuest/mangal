package metadata

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata"
)

var _ tea.Model = (*Model)(nil)

// Model implements tea.Model.
type Model struct {
	meta   metadata.Metadata
	styles styles

	currentStyle style
	ShowFull     bool
}

// SetMetadata replaces the current metadata and updates the style.
func (m *Model) SetMetadata(meta metadata.Metadata) {
	m.meta = meta
	m.updateStyle()
}

// Metadata returns the current metadata.
func (m *Model) Metadata() metadata.Metadata {
	return m.meta
}

// Style returns the style based on the current metadata.
func (m *Model) Style() style {
	return m.currentStyle
}

func (m *Model) IDStyle(id metadata.IDSource) style {
	switch id {
	case metadata.IDSourceAnilist:
		return m.styles.anilist
	case metadata.IDSourceMyAnimeList:
		return m.styles.myAnimeList
	case metadata.IDSourceKitsu:
		return m.styles.kitsu
	case metadata.IDSourceMangaUpdates:
		return m.styles.mangaUpdates
	case metadata.IDSourceAnimePlanet:
		return m.styles.animePlanet
	default:
		s := m.styles.provider
		s.Code = m.meta.ID().Code
		return s
	}
}

// updateStyle sets the style based on the current metadata.
func (m *Model) updateStyle() {
	m.currentStyle = m.IDStyle(m.meta.ID().Source)
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
		return m.styles.base.Background(m.currentStyle.Color).Render(m.meta.Title() + year)
	}

	sep := ": "
	if id := m.meta.ID().Raw; id != "" {
		sep = " (" + id + "): "
	}
	return m.styles.base.Background(m.currentStyle.Color).Render(m.currentStyle.Prefix + sep + m.meta.Title() + year)
}
