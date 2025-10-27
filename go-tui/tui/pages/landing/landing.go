package landing

import (
	"tui/tui/constants"
	"tui/tui/constants/custommenu"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/kevm/bubbleo/menu"
	// "github.com/kevm/bubbleo/menu"
	// "github.com/kevm/bubbleo/menu"
)

type Model struct {
	// mode mode
	// stats
	// viewport   viewport.Model
	// headerText string
	menu custommenu.Model
	// footerText string
	// currentModel string
	// selectedPage
	// pages []Page
	// SelectedPage tea.Model
	// selectedPage
	// quitting bool
}

// type model struct {
// 	SelectedPage tea.Model
// }

// func createShell(m tea.Model) tea.Model {
// 	s := shell.New()
// 	s.Navstack.Push(navstack.NavigationItem{Model: m, Title: "Title"})
// 	p := tea.NewProgram(s, tea.WithAltScreen())

// 	finalshell, err := p.Run()
// 	if err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}

// 	topNavItem := finalshell.(shell.Model).Navstack.Top()
// 	if topNavItem == nil {
// 		log.Printf("Nothing selected")
// 		os.Exit(1)
// 	}

// 	selected := topNavItem.Model.(Model)
// 	return selected
// } // createShell

func New(options []tea.Model) Model {
	choices := []custommenu.Choice{
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
	title := ""
	// header := ""
	// footer := "enter: select\t|\tesc: back\n"
	// var currChoice *menu.Choice
	menuModel := custommenu.New(title, choices, nil)
	landingModel := Model{
		// viewport:   viewport.New(0, 0),
		// headerText: header,
		menu: menuModel,
		// footerText: footer,
	}
	return landingModel
} // NewMenu

func (m Model) Init() tea.Cmd {
	return nil
} // Init

// func selected
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case constants.PageSelectedMsg:
		selected := msg
		if selected.Model != nil {
			return m, constants.PushNavStackCmd(selected)
		} // if
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		} // switch
	case tea.WindowSizeMsg:
		m.menu.SetSize(msg)
		return m, nil
	} // switch
	updatedMenu, cmd := m.menu.Update(msg)
	m.menu = updatedMenu.(custommenu.Model)
	return m, cmd
}

func (m Model) View() string {
	// styles := constants.NewLandingStyles(m.viewport.Width)
	// // s := styles.Header
	// // s += constants.MenuStyle.Render(m.menu.View())
	// // s += styles.Footer
	// // header := styles.Header.Render("Claim Extraction Pipeline")
	// // menu := styles.Menu.Render(m.menu.View())
	// // footer := styles.Footer.Render("help section...")

	// // view := lipgloss.JoinVertical(
	// // 	// lipgloss.Top,
	// // 	lipgloss.Center,
	// // 	styles.Header.Render(m.headerText),
	// // 	styles.Menu.Render(m.menu.View()),
	// // 	styles.Footer.Render(m.footerText),
	// // )

	// return lipgloss.Place(
	// 	m.viewport.Width,
	// 	m.viewport.Height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Top,
	// 		styles.Header.Render(m.headerText),
	// 		styles.Menu.Render(m.menu.View()),
	// 		styles.Footer.Render(m.footerText),
	// 	),
	// )
	return m.menu.View()
	// return view
	// s += "enter: select\t|\tesc: back\n"
	// return s
}
