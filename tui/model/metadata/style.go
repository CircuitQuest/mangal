package metadata

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type styles struct {
	provider,
	anilist lipgloss.Style
}

func defaultStyles() styles {
	baseStyle := style.Normal.Base.Padding(0, 1).Foreground(color.Bright)
	return styles{
		provider: baseStyle.Background(color.Provider),
		anilist:  baseStyle.Background(color.Anilist),
	}
}
