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
	Bright = lipgloss.Color("#FEFEFE")

	// Metadata provider colors
	Provider     = lipgloss.Color("#f26f63")
	Anilist      = lipgloss.Color("#02A9FF") // #1E2630
	MyAnimeList  = lipgloss.Color("#2E51A2")
	Kitsu        = lipgloss.Color("#F75239") // #312631
	MangaUpdates = lipgloss.Color("#F28A2E") // no clear code available
	AnimePlanet  = lipgloss.Color("#A72A2D") // #1C3867 #2E4F83 #F69330 #E65448 #A72A2D // not exact
)
