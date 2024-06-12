package base

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	Resize(size Size)

	// Model (wrapper of tea.Model) methods.
	Init(model Model) tea.Cmd
	Update(model Model, msg tea.Msg) tea.Cmd
	View(model Model) string
}

type Title struct {
	Text                   string
	Background, Foreground lipgloss.Color
}

type Size struct {
	Width, Height int
}
