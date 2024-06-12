package base

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is a wrapper interface of tea.Model.
type Model interface {
	tea.Model

	Context() context.Context
}
