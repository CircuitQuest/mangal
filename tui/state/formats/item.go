package formats

import (
	"strings"

	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/theme/style"
)

type Item struct {
	format libmangal.Format
}

func (i Item) FilterValue() string {
	return i.format.String()
}

func (i Item) Title() string {
	var sb strings.Builder

	sb.WriteString(i.FilterValue())

	if i.IsSelectedForDownloading() {
		sb.WriteString(" ")
		sb.WriteString(style.Bold.Accent.Render("Download"))
	}

	if i.IsSelectedForReading() {
		sb.WriteString(" ")
		sb.WriteString(style.Bold.Accent.Render("Read"))
	}

	return sb.String()
}

func (i Item) Description() string {
	ext := i.format.Extension()

	if ext == "" {
		return "<none>"
	}

	return ext
}

func (i Item) IsSelectedForDownloading() bool {
	format := config.Download.Format.Get()

	return i.format == format
}

func (i Item) IsSelectedForReading() bool {
	format := config.Read.Format.Get()

	return i.format == format
}

func (i Item) SelectForDownloading() error {
	if err := config.Download.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}

func (i Item) SelectForReading() error {
	if err := config.Read.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}
