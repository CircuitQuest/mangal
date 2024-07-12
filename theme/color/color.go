package color

import "github.com/charmbracelet/lipgloss"

var (
	// Available as styles, too
	Accent     = lipgloss.Color("#EB5E28")
	Secondary  = lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"}
	Background = lipgloss.Color("#252422")
	Success    = lipgloss.Color("#7EC699")
	Warning    = lipgloss.Color("#EBCA89")
	Error      = lipgloss.Color("#E05252")
	Loading    = lipgloss.Color("#A49FA5")
	Viewport   = lipgloss.Color("#008080")

	// Available only as colors
	Bright   = lipgloss.Color("#FEFEFE")
	Provider = lipgloss.Color("#5604b5")
	Anilist  = lipgloss.Color("#02A9FF")
)
