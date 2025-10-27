package tui

import (
	"fmt"
	"os"
	"tui/tui/constants"
	"tui/tui/constants/customshell"
	"tui/tui/pages/filedata"
	"tui/tui/pages/filespans"
	"tui/tui/pages/landing"
	"tui/tui/pages/pullfiles"
	"tui/tui/pages/runpipeline"
	"tui/tui/pages/searchfiles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevm/bubbleo/navstack"
)

// type sessionState int

// const (
// 	landingView sessionState = iota
// 	pullFilesView
// 	runModelView
// 	searchFilesView
// 	fileDataView
// 	fileSpansView
// )

// type MainModel struct {
// 	state       sessionState
// 	landing     tea.Model
// 	pullFiles   tea.Model
// 	runPipeline tea.Model
// 	searchFiles tea.Model
// 	fileData    tea.Model
// 	fileSpans   tea.Model
// }

//	func initializeMainModel() MainModel {
//		var menuModels = []tea.Model{
//			pullfiles.New(),
//			searchfiles.New(),
//			runpipeline.New(),
//		}
//		return MainModel{
//			state:     landingView,
//			landing:   landing.New(menuModels),
//			fileData:  filedata.New(),
//			fileSpans: filespans.New(),
//		}
//	}
// type CustomShell struct {
// 	*shell.Model
// }

// func (m *CustomShell) View() string {
// 	m.Breadcrumb.Styles.Delimiter = " -> "
// 	bc := m.Breadcrumb.View()
// 	nav := m.Navstack.View()
// 	return lipgloss.NewStyle().Render(bc, nav)
// }

func initializeMainShell() tea.Model {
	// create page models
	pullFilesPage := pullfiles.New()
	runPipelinePage := runpipeline.New()
	searchFilesPage := searchfiles.New()
	fileDataPage := filedata.New()
	fileSpansPage := filespans.New()

	// landing menu choices pointing to other pages
	landingMenu := landing.New([]tea.Model{
		pullFilesPage,
		runPipelinePage,
		searchFilesPage,
	})

	// shell and navstack
	s := customshell.New()
	// cs := CustomShell{s}
	// cs := constants.CustomShell{Model: s}
	// push landing menu as initial page
	s.Navstack.Push(navstack.NavigationItem{
		Model: landingMenu,
		Title: "",
	})

	// searchfiles needs to know which page to push for selected file
	// for example, it can have a callback or send a msg to push filedata
	// searchFilesPage.SetFileDataPage(fileDataPage)
	_ = fileSpansPage
	_ = fileDataPage

	return s
} // initializeMainShell

// StartTea the entry point for the UI
func StartTea() error {
	constants.P = tea.NewProgram(initializeMainShell(), tea.WithAltScreen())
	_, err := constants.P.Run()
	if err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	} // if
	return nil
} // StartTea
