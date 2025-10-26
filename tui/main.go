package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	// "github.com/charmbracelet/lipgloss"
	"fmt"
)

type model struct {
	viewport    viewport.Model
	currentView viewState
	landing     LandingModel
	// runPipeline RunPipelineModel
}

type viewState int

const (
	landing viewState = iota
	runPipeline
	processedFiles
	logs
)

type LandingModel struct {
	banner string // welcome message

}

// load landing page
func initialModel() model {
	return model{
		viewport:    viewport.New(width, height),
		currentView: landing,
		landing: LandingModel{
			banner: "Claim Extraction Pipeline",
		},
	}
} // initialModel

// GOAL: display landing page with this
// - banner
// - stats
//   - num of raw files
//   - num of processed files
//   - num of files pending extraction
//
// - options
//   - run pipeline on new files
//   - view processed files
//   - quit
func (m model) Init() tea.Cmd {
	return nil
} // Init

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
} // Update

func (m model) View() string {
	var b strings.Builder
	// var s string

	return b.String()
} // View

func main() {
	p := tea.NewProgram(initialModel())

	_, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} // if

} // main
