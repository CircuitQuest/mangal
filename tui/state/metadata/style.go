package metadata

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/tui/util"
)

type styles struct {
	meta util.MetaStyle
	enum func(list.Items, int) string

	fieldName,
	enumerator lipgloss.Style
}

func defaultStyles(meta util.MetaStyle) styles {
	base := lipgloss.NewStyle().Foreground(meta.Color)
	return styles{
		meta: meta,
		enum: func(_ list.Items, _ int) string {
			return icon.SubItem.Raw()
		},
		fieldName:  base,
		enumerator: base.MarginRight(1).PaddingLeft(2),
	}
}
