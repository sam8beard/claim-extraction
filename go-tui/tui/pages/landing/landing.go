package landing

import (
	"fmt"
	"log"
	"os"
	"tui/tui/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevm/bubbleo/menu"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/shell"
)

type Model struct {
	// mode mode
	menu menu.Model
	// pages []Page
	// SelectedPage tea.Model
	// selectedPage
	// quitting bool
}

// func NewChoices([]Model) []menu.Choice {
//
//		choices = []menu.Choice{
//		{
//			Title:       "Pull Files",
//			Description: `Search for and pull files to process`,
//			Model:       pullfiles.New(),
//		},
//		{
//			Title:       "Run Pipeline",
//			Description: `Search for and pull files to process`,
//			Model:       runpipeline.New(),
//		},
//		{
//			Title:       "Search Files",
//			Description: `Search for and pull files to process`,
//			Model:       searchfiles.New(),
//		},
//	}

func NewMenu(options []tea.Model) menu.Model {
	var choices = []menu.Choice{
		{
			Title:       "Pull Files",
			Description: `placeholder`,
			Model:       options[0],
		},
		{
			Title:       "Run Pipeline",
			Description: `placeholder`,
			Model:       options[1],
		},
		{
			Title:       "Search Files",
			Description: `placeholder`,
			Model:       options[2],
		},
	}
	title := "Options"
	m := menu.New(title, choices, nil)
	return m
}

func New(landingOptions []tea.Model) Model {
	menu := NewMenu(landingOptions)
	return Model{
		menu: menu,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func createShell(m tea.Model) tea.Model {
	s := shell.New()
	s.Navstack.Push(navstack.NavigationItem{Model: m, Title: ""})
	p := tea.NewProgram(s, tea.WithAltScreen())

	finalshell, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	topNavItem := finalshell.(shell.Model).Navstack.Top()
	if topNavItem == nil {
		log.Printf("Nothing selected")
		os.Exit(1)
	}

	selected := topNavItem.Model.(Model)
	return selected
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			// this causes a loop
			newM := createShell(m)
			return newM, nil
		} // switch
	case tea.WindowSizeMsg:
		m.menu.SetSize(msg)
		return m, nil
	} // switch
	updatedMenu, cmd := m.menu.Update(msg)
	m.menu = updatedMenu.(menu.Model)
	return m, cmd
}

func (m Model) View() string {
	return constants.DocStyle.Render(m.menu.View())
}
