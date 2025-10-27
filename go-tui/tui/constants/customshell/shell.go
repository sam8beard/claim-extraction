// Package Shell is a basic wrapper around the navstack and breadcrumb packages
// It provides a basic navigation mechanism while showing breadcrumb view of where the user is
// within the navigation stack.
package customshell

import (
	"tui/tui/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/breadcrumb"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/utils"
	"github.com/kevm/bubbleo/window"
)

type Model struct {
	Navstack   *navstack.Model
	Breadcrumb breadcrumb.Model
	window     *window.Model

	Header string
	Footer string
}

// New creates a new shell model
func New() Model {
	w := window.New(120, 30, 0, 0)
	ns := navstack.New(&w)
	bc := breadcrumb.New(&ns)

	return Model{
		Navstack:   &ns,
		Breadcrumb: bc,
		window:     &w,
		Header:     "Claim Extraction Pipeline",
		Footer:     "enter: select\t|\tesc: back\n",
	}
}

// Init determines the size of the widow used by the navigation stack.
func (m Model) Init() tea.Cmd {

	w, h := m.Breadcrumb.Styles.Frame.GetFrameSize()
	m.window.SideOffset = w
	m.window.TopOffset = h

	return utils.Cmdize(m.window.GetWindowSizeMsg())
}

// Update passes messages to the navigation stack.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	m.window.Width = msg.Width
	// 	m.window.Height = msg.Height
	// }
	cmd := m.Navstack.Update(msg)
	return m, cmd
}

// View renders the breadcrumb and the navigation stack.
func (m Model) View() string {
	styles := constants.NewAppStyles()
	m.Breadcrumb.Styles.Delimiter = " / "
	// bc := m.Breadcrumb.View()
	// body := m.Navstack.View() // or any page
	// content := lipgloss.JoinVertical(lipgloss.Top, body)
	body := styles.Body.Render(m.Navstack.View())
	bc := styles.Breadcrumb.Render(m.Breadcrumb.View())
	header := styles.Header.Render(m.Header)
	footer := styles.Footer.Render(m.Footer)
	// body := styles.Body.Render(lipgloss.JoinVertical(lipgloss.Top, bc, nav))

	return lipgloss.Place(
		m.window.Width,
		m.window.Height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			bc,
			body,
			footer,
		),
	)
}
