package volumes

import (
	"fmt"
	"github.com/luevano/libmangal"
)

type Item struct {
	volume *libmangal.Volume
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
