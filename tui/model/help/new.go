package help

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
)

func New() help.Model {
	h := help.New()
	h.Ellipsis = icon.Ellipsis.Raw()
	h.ShortSeparator = " " + icon.Separator.Raw() + " "
	h.FullSeparator = " " + icon.Separator.Raw() + " "

	keyStyle := style.Bold.Warning
	sepStyle := style.Normal.Base
	descStyle := style.Normal.Secondary

	h.Styles = help.Styles{
		Ellipsis:       style.Normal.Base,
		ShortKey:       keyStyle,
		ShortDesc:      descStyle,
		ShortSeparator: sepStyle,
		FullKey:        keyStyle,
		FullDesc:       descStyle,
		FullSeparator:  sepStyle,
	}
	return h
}
