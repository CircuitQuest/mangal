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
		header.WriteString(m.styles.notification.Width(width).Render(m.notification))
	}

	if subtitle := m.state.Subtitle(); subtitle != "" {
		header.WriteString(newline)
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
	var footer strings.Builder
	footer.Grow(200)

	if m.loadingMessage != "" {
		footer.WriteString(m.styles.spinner.Render(m.spinner.View()))
		footer.WriteString(m.styles.loading.Render(m.loadingMessage))
		footer.WriteString(newline)
	}
	footer.WriteString(m.help.View(m.keyMap.with(m.state.KeyMap())))

	return m.styles.footer.Render(footer.String())
}
