package pathtable

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/path"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/tui/util"
)

func Run() error {
	model := newModel()
	_, err := tea.NewProgram(model).Run()
	return err
}

func newModel() *Model {
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

	rows := []table.Row{
		{
			"config",
			path.ConfigDir(),
		},
		{
			"providers",
			path.ProvidersDir(),
		},
		{
			"downloads",
			path.DownloadsDir(),
		},
		{
			"cache",
			path.CacheDir(),
		},
		{
			"logs",
			path.LogDir(),
		},
		{
			"temp",
			path.TempDir(),
		},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)),
		table.WithStyles(table.Styles{
			Header:   style.Bold.Accent,
			Cell:     style.Normal.Base,
			Selected: style.Italic.Accent,
		}),
	)

	return &Model{
		table: t,
		help:  help.New(),
		keyMap: keyMap{
			Copy: util.Bind("copy", "enter"),
			Quit: util.Bind("quit", "q", "ctrl+c"),
		},
	}
}
