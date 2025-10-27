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
		Footer:     "enter: select\tesc: back\n",
	}
}

// Init determines the size of the widow used by the navigation stack.
func (m Model) Init() tea.Cmd {
	// m.window.Height = lipgloss.Height(m.Navstack.View())
	// m.window.Height = lipgloss.Height(m.Navstack.View())

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
	if m.window.Width == 0 {
		return "Initializing..."
	} // if
	styles := constants.NewAppStyles()
	m.Breadcrumb.Styles.Delimiter = " / "
	// bc := m.Breadcrumb.View()
	// body := m.Navstack.View() // or any page
	// content := lipgloss.JoinVertical(lipgloss.Top, body)
	// content := styles.Content.Render(m.Navstack.View())
	bc := styles.Breadcrumb.Render(m.Breadcrumb.View())
	title := styles.Header.Render(m.Header)
	footer := styles.Footer.Render(m.Footer)
	// body := styles.Body.Render(lipgloss.JoinVertical(lipgloss.Top, bc, nav))
	// content := lipgloss.Place(
	// 	m.window.Width,
	// 	m.window.Height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	body,
	// )

	// contains the title and the breadcrumb
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		title,
		bc,
	)

	// contentHeight := m.window.Height - lipgloss.Height(m.Footer)
	// styles.SetContentHeight(contentHeight)
	// styles.Content.Height(contentHeight)
	contentHeight := m.window.Height - lipgloss.Height(m.Navstack.View())
	// content := styles.Content.Height(0).Render(m.Navstack.View())
	content := styles.Content.Render(m.Navstack.View())
	// footer := styles.Footer.AlignVertical(lipgloss.Bottom).Render(m.Footer)
	// cheight := fmt.Sprinf("%d", )
	// windowHeight := fmt.Sprintf("%d", m.window.Height)
	// footerHeight := fmt.Sprintf("%d", lipgloss.Height(footer))
	// m.window.Width = lipgloss.Width(m.Navstack.View())
	// m.window.Height = lipgloss.Height(m.Navstack.View())
	// if m.window.Height <
	// s := header + " "
	// lipgloss.PlaceHorizontal()
	header = lipgloss.Place(
		m.window.Width,
		lipgloss.Height(header),
		// 10,
		lipgloss.Center,
		lipgloss.Center,
		header,
	)
	footer = lipgloss.Place(
		m.window.Width,
		lipgloss.Height(footer),
		// 10,
		lipgloss.Center,
		lipgloss.Center,
		footer,
	)

	content = lipgloss.Place(
		m.window.Width,
		lipgloss.Height(content),
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
	contentPadding := m.window.Height - contentHeight
	content = styles.Content.PaddingBottom(contentPadding).Render(content)
	footer = styles.Footer.PaddingBottom(0).Render(footer)
	// footer = lipgloss.Place()
	// s := lipgloss.Place(
	// 	m.window.Width,
	// 	m.window.Height,
	// 	lipgloss.Center,
	// 	lipgloss.Center,
	// 	lipgloss.JoinVertical(
	// 		lipgloss.Center,
	// 		header,
	// 		content,
	// 		footer,
	// 		fmt.Sprintf("%d", contentHeight),
	// 		fmt.Sprintf("%d", m.window.Height),
	// 		fmt.Sprintf("%d", lipgloss.Height(header)),
	// 		fmt.Sprintf("%d", lipgloss.Height(m.Navstack.View())),
	// 		fmt.Sprintf("%d", m.window.TopOffset),
	// 		// footerHeight,
	// 		// windowHeight,
	// 		// contentHeight,
	// 	),
	// )
	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)

	available := m.window.Height - (headerH + footerH)
	padTop := max((available-contentHeight)/2, 0)
	padBottom := max(available-padTop-contentHeight, 0)

	content = lipgloss.NewStyle().
		PaddingTop(padTop).
		PaddingBottom(padBottom).
		Render(content)

	s := lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		content,
		footer,
		// fmt.Sprintf("%d", contentHeight),
		// fmt.Sprintf("%d", m.window.Height),
		// fmt.Sprintf("%d", lipgloss.Height(header)),
		// fmt.Sprintf("%d", lipgloss.Height(m.Navstack.View())),
		// fmt.Sprintf("%d", m.window.TopOffset),
	)

	// s = lipgloss.PlaceVertical()
	return s
}
