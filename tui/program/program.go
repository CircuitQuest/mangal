package program

import tea "github.com/charmbracelet/bubbletea"

var tui *tea.Program

// TUI get the TUI program.
func TUI() *tea.Program {
	return tui
}

// SetTUI sets the TUI program.
func SetTUI(program *tea.Program) {
	tui = program
}

// Send sends a message to the TUI program.
//
// Send is intended to be used externally to send messages
// to the TUI program but it's also useful when there is no
// other easy way to return a message/cmd (like notifications).
func Send(msg tea.Msg) {
	tui.Send(msg)
}
