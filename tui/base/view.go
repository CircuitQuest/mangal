package base

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const newline = "\n"

// View implements tea.Model.
func (m *model) View() string {
	header := m.viewHeader()
	state := m.viewState()
	footer := m.viewFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, state, footer)
}

func (m *model) viewHeader() string {
	var header strings.Builder
	header.Grow(200)

	title := m.state.Title()
	titleStyle := m.styles.title

	if title.Background != "" {
		titleStyle = titleStyle.Background(title.Background)
	}
	if title.Foreground != "" {
		titleStyle = titleStyle.Foreground(title.Foreground)
	}
	header.WriteString(titleStyle.MaxWidth(m.size.Width / 2).Render(title.Text))

	if status := m.state.Status(); status != "" {
		header.WriteString(m.styles.status.Render(status))
	}

	if m.notification != "" {
		width := m.size.Width - lipgloss.Width(header.String())
		header.WriteString(m.styles.notification.MaxWidth(width).Render(m.notification))
	}

	header.WriteString(newline)
	if m.loadingMessage != "" {
		header.WriteString(m.spinner.View())
	}
	header.WriteString(m.styles.loading.Render(m.loadingMessage))

	header.WriteString(newline)
	if subtitle := m.state.Subtitle(); subtitle != "" {
		header.WriteString(m.styles.subtitle.Render(subtitle))
	}

	return m.styles.header.Render(header.String())
}

func (m *model) viewState() string {
	size := m.stateSize()
	return lipgloss.Place(
		size.Width,
		size.Height,
		lipgloss.Left,
		lipgloss.Top,
		m.styles.state.Render(m.state.View()),
	)
}

func (m *model) viewFooter() string {
	return m.styles.footer.Render(m.help.View(m.keyMap.with(m.state.KeyMap())))
}
