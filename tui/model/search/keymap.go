package search

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/luevano/mangal/tui/util"
)

func newKeyMap() keyMap {
	return keyMap{
		confirm: util.Bind("confirm", "enter"),
		cancel:  util.Bind("cancel", "esc"),
	}
}

type keyMap struct {
	confirm,
	cancel key.Binding
}

// enableKeyMap if the keymap should be enabled.
func (m *Model) enableKeyMap(enable bool) {
	if enable {
		m.keyMap.cancel.SetEnabled(true)
		m.keyMap.confirm.SetEnabled(m.input.Value() != "")
	} else {
		m.keyMap.cancel.SetEnabled(false)
		m.keyMap.confirm.SetEnabled(false)
	}
}
