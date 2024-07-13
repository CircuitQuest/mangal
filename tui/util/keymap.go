package util

import "github.com/charmbracelet/bubbles/key"

// Bind is a convenience function to create key bindings easily.
func Bind(help string, primaryKey string, extraKeys ...string) key.Binding {
	keys := make([]string, 1+len(extraKeys))
	keys[0] = primaryKey
	for i, k := range extraKeys {
		keys[i+1] = k
	}

	return BindNamedKey(primaryKey, help, keys...)
}

// BindNamedKey creates a keybind with specified help and list of keys.
func BindNamedKey(keyHelp, help string, keys ...string) key.Binding {
	if keyHelp == " " {
		keyHelp = "space"
	}
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keyHelp, help),
	)
}
