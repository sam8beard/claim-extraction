package constants

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/utils"
)

var (
	// P the program
	P *tea.Program

	// WindowSize the size of the terminal window
	WindowSize tea.WindowSizeMsg
)

// Styling
type AppStyles struct {
	Header     lipgloss.Style
	Breadcrumb lipgloss.Style
	Footer     lipgloss.Style
	Body       lipgloss.Style
}

func NewAppStyles() AppStyles {
	return AppStyles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("24")),
		Breadcrumb: lipgloss.NewStyle().
			Italic(true),
		Footer: lipgloss.NewStyle().
			Italic(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("240")),
		Body: lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(2, 2),
	}
}

// type LandingStyles struct {
// 	Viewport viewport.Model
// 	Footer   lipgloss.Style
// 	Menu     lipgloss.Style
// 	Header   lipgloss.Style
// }

// func NewLandingStyles(width int) LandingStyles {
// 	// return Styles{
// 	// 	Footer: lipgloss.NewStyle().
// 	// 		Width(width).
// 	// 		AlignHorizontal(lipgloss.Center),
// 	// 	Menu: lipgloss.NewStyle().Margin(0, 2),
// 	// 	Header: lipgloss.NewStyle().
// 	// 		Bold(true),
// 	// }
// 	view := viewport.New(1, 1)
// 	var headerStyles = lipgloss.NewStyle().
// 		Width(width).
// 		AlignHorizontal(lipgloss.Center).
// 		Bold(true)

// 	var footerStyles = lipgloss.NewStyle().
// 		Width(width).
// 		AlignHorizontal(lipgloss.Center).
// 		Foreground(lipgloss.Color("240"))

// 	var headerHeight = headerStyles.GetHeight()
// 	var footerHeight = footerStyles.GetHeight()

// 	var menuStyles = lipgloss.NewStyle().
// 		Height(view.Height - headerHeight - footerHeight).
// 		PaddingTop(2)

// 	return LandingStyles{
// 		Viewport: view,
// 		Header:   headerStyles,
// 		Menu:     menuStyles,
// 		Footer:   footerStyles,
// 	}
// }

// Key bindings
type keymap struct {
	// put custom keys here
}

var KeyMap = keymap{
	// define custom keys here
}

// represents a page selection in landing menu
type PageSelectedMsg struct {
	Model tea.Model
	Name  string
}

// represents a cmd that pushes the page that was selected
// onto the navstack
func PushNavStackCmd(msg PageSelectedMsg) tea.Cmd {
	return utils.Cmdize(navstack.PushNavigation{
		Item: navstack.NavigationItem{
			Title: msg.Name,
			Model: msg.Model,
		},
	})
} // PushNavStackCmd
