package formats

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/theme/style"
)

var (
	_ list.Item        = (*item)(nil)
	_ list.DefaultItem = (*item)(nil)
)

// item implements list.item.
type item struct {
	format libmangal.Format
}

// FilterValue implements list.Item.
func (i *item) FilterValue() string {
	return i.format.String()
}

// Title implements list.DefaultItem.
func (i *item) Title() string {
	var sb strings.Builder

	sb.WriteString(i.FilterValue())

	if i.isSelectedForDownloading() {
		sb.WriteString(" ")
		sb.WriteString(style.Bold.Accent.Render("Download"))
	}

	if i.isSelectedForReading() {
		sb.WriteString(" ")
		sb.WriteString(style.Bold.Accent.Render("Read"))
	}

	return sb.String()
}

// Description implements list.DefaultItem.
func (i *item) Description() string {
	ext := i.format.Extension()

	if ext == "" {
		return "<none>"
	}

	return ext
}

func (i *item) isSelectedForDownloading() bool {
	format := config.Download.Format.Get()

	return i.format == format
}

func (i *item) isSelectedForReading() bool {
	format := config.Read.Format.Get()

	return i.format == format
}

// TODO: don't set? or just don't write config.
func (i *item) selectForDownloading() error {
	if err := config.Download.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}

// TODO: don't set? or just don't write config.
func (i *item) selectForReading() error {
	if err := config.Read.Format.Set(i.format); err != nil {
		return err
	}

	return config.Write()
}
