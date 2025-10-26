package tui

import (
	"fmt"
	"os"
	"tui/tui/constants"
	"tui/tui/pages/filedata"
	"tui/tui/pages/filespans"
	"tui/tui/pages/landing"
	"tui/tui/pages/pullfiles"
	"tui/tui/pages/runpipeline"
	"tui/tui/pages/searchfiles"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	landingView sessionState = iota
	pullFilesView
	runModelView
	searchFilesView
	fileDataView
	fileSpansView
)

type MainModel struct {
	state       sessionState
	landing     tea.Model
	pullFiles   tea.Model
	runPipeline tea.Model
	searchFiles tea.Model
	fileData    tea.Model
	fileSpans   tea.Model
}

func initializeMainModel() MainModel {
	var menuModels = []tea.Model{
		pullfiles.New(),
		searchfiles.New(),
		runpipeline.New(),
	}
	return MainModel{
		state:     landingView,
		landing:   landing.New(menuModels),
		fileData:  filedata.New(),
		fileSpans: filespans.New(),
	}
}

// StartTea the entry point for the UI
func StartTea() error {
	constants.P = tea.NewProgram(initializeMainModel(), tea.WithAltScreen())
	_, err := constants.P.Run()
	if err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	} // if
	return nil
} // StartTea

// Init initialize the main model
func (m MainModel) Init() tea.Cmd {
	return nil
} // Init

// Update handles the updates for the main model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.state {
	case landingView:
		newLanding, newCmd := m.landing.Update(msg)
		m.landing = newLanding
		cmd = newCmd

	} // switch
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
} // Update

// View view the main model
func (m MainModel) View() string {
	switch m.state {
	case landingView:
		return m.landing.View()
	} // switch

	return ""
} // View
