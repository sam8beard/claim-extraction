package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// P the program
	P *tea.Program

	// WindowSize the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

// Styling
var DocStyle = lipgloss.NewStyle().Margin(0, 2)

// Key bindings
type keymap struct {
	// put custom keys here
}

var KeyMap = keymap{
	// define custom keys here
}
