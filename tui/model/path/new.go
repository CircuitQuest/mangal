package path

import (
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/model/help"
)

func New(standalone bool) *Model {
	cols := []table.Column{
		{
			Title: "Name",
			Width: 10,
		},
		{
			Title: "Path",
			Width: 40,
		},
	}

	paths := path.AllPaths()
	rows := make([]table.Row, len(paths))
	keys := paths.Keys()
	slices.Sort(keys)
	for i, k := range keys {
		rows[i] = table.Row{string(k), paths.Get(k)}
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
		table.WithStyles(table.Styles{
			Header:   style.Bold.Accent,
			Cell:     style.Normal.Base,
			Selected: style.Normal.Background.Background(color.Accent),
		}),
	)

	m := &Model{
		table:                t,
		help:                 help.New(),
		standalone:           standalone,
		notificationDuration: time.Second,
		style:                style.Normal.Base.Margin(0, 2),
		msgStyle:             style.Normal.Warning,
		keyMap:               newKeyMap(),
	}
	if !standalone {
		m.DisableQuitKeybindings()
	}
	return m
}
