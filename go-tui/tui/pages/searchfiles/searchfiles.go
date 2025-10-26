package searchfiles

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	// mode mode
	placeholder string
	// pages []Page
	// SelectedPage tea.Model
	// selectedPage
	// quitting bool
}

func New() Model {
	return Model{
		placeholder: "testing",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}
