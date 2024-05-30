package volumes

import (
	"fmt"

	"github.com/luevano/libmangal/mangadata"
)

type Item struct {
	volume *mangadata.Volume
}

func (i Item) FilterValue() string {
	return fmt.Sprintf("Volume %s", (*i.volume).String())
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return ""
}
