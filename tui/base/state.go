package base

import (
	"context"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// State is kind of an indirect wrapper interface of Model (wrapper of tea.Model).
//
// State could be thought of as a "Window".
type State interface {
	Intermediate() bool
	Backable() bool

	KeyMap() help.KeyMap

	Title() Title
	Subtitle() string
	Status() string

	// Resize the State usable viewport. Gets called at least once before Init.
	Resize(size Size)

	// Model (wrapper of tea.Model) methods.
	Init(ctx context.Context) tea.Cmd
	Update(ctx context.Context, msg tea.Msg) tea.Cmd
	View() string
}
