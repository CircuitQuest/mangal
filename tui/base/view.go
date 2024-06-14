package base

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	stringutil "github.com/luevano/mangal/util/string"
	"github.com/muesli/reflow/wordwrap"
)

// View implements tea.Model.
func (m *Model) View() string {
	const newline = "\n"

	title := m.state.Title()
	titleStyle := m.styles.Title

	if title.Background != "" {
		titleStyle = titleStyle.Background(title.Background)
	}

	if title.Foreground != "" {
		titleStyle = titleStyle.Foreground(title.Foreground)
	}

	titleText := stringutil.Trim(title.Text, m.size.Width/2)
	header := m.styles.TitleBar.Render(titleStyle.Render(titleText) + " " + m.state.Status())

	subtitle := m.state.Subtitle()
	if subtitle != "" {
		header = lipgloss.JoinVertical(lipgloss.Left, header, m.styles.TitleBar.Render(m.styles.Subtitle.Render(m.state.Subtitle())))
		// header += m.styles.TitleBar.Render(m.styles.Subtitle.Render(m.state.Subtitle()))
	}

	view := wordwrap.String(m.state.View(), m.size.Width)
	keyMapHelp := m.styles.HelpBar.Render(m.help.View(m))

	headerHeight := lipgloss.Height(header)
	viewHeight := lipgloss.Height(view)
	helpHeight := lipgloss.Height(keyMapHelp)

	diff := m.size.Height - headerHeight - viewHeight - helpHeight

	var filler string
	if diff > 0 {
		filler = strings.Repeat(newline, diff)
	}

	// return lipgloss.JoinVertical(lipgloss.Left, header, view, filler, keyMapHelp)
	return header + newline + view + filler + newline + keyMapHelp
}
