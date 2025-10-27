package pullfiles

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
} // Init

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	switch msg.String() {
	// 	case "enter":
	// 		pop := utils.Cmdize(navstack.PopNavigation{})
	// 		selected := utils.Cmdize()
	// 	} // switch
	// } // switch
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		} // if
	}
	return m, nil
} // Update

func (m Model) View() string {
	return "FIRING"
} // View
