package pullfiles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/utils"
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

			// if esc, go back a page
		case "esc":
			cmd = utils.Cmdize(navstack.PopNavigation{})
			return m, cmd
		}
	}
	return m, nil
} // Update

func (m Model) View() string {
	return "Place holder"
} // View
