package metadata

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
)

func (m *Model) ShowMetadataCmd() tea.Cmd {
	return base.Viewport(m.Title(), m.RenderMetadata(), m.Color())
}
