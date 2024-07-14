package chapters

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/style"
)

type styles struct {
	sep,
	subtitle,
	format lipgloss.Style
}

func defaultStyles() styles {
	return styles{
		sep:      style.Bold.Warning.Padding(0, 1),
		subtitle: style.Normal.Secondary, // matches base without padding
		format:   style.Bold.Warning,
	}
}
