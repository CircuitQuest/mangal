package volumes

import (
	"fmt"
	"github.com/luevano/libmangal"
)

type Item struct {
	libmangal.Volume
}

func (i Item) FilterValue() string {
	return fmt.Sprintf("Volume %.1f", i.Info().Number)
}

func (i Item) Title() string {
	return i.FilterValue()
}

func (i Item) Description() string {
	return ""
}
