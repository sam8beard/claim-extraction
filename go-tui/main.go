package main

import (
	"fmt"
	"os"
	"tui/tui"
)

func main() {
	err := tui.StartTea()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} // if
} // main

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
// 	runModel    tea.Model
// 	searchFiles tea.Model
// 	fileData    tea.Model
// 	fileSpans   tea.Model
// }

// var docStyle = lipgloss.NewStyle().Margin(1, 2)

// type LandingModel struct {
// 	stats   []Stat     // stats of files
// 	options list.Model // list of options
// }

// // stat on landing page
// type Stat struct {
// 	description string
// 	numFiles    int
// }

// // create new stat
// func NewStat(description string, numFiles int) Stat {
// 	return Stat{
// 		description: description,
// 		numFiles:    numFiles,
// 	}
// }

// // option on options list
// type option struct {
// 	title string
// 	desc  string
// }

// func (o option) Title() string       { return o.title }
// func (o option) Description() string { return o.desc }
// func (o option) FilterValue() string { return o.title }

// // create new option
// func NewOption(title string, desc string) list.Item {
// 	return option{title: title, desc: desc}
// }

// func NewOptionList(options []list.Item) list.Model {
// 	optionList := list.New(options, list.NewDefaultDelegate(), 0, 0)
// 	optionList.Title = "Select an option:"
// 	return optionList
// } // buildOptionList

// var options = []list.Item{
// 	NewOption(
// 		"Run pipeline",
// 		"Run trained pipeline on newly pulled files",
// 	),
// 	NewOption(
// 		"View processed files",
// 		"View files that have already been processed by the pipeline",
// 	),
// 	NewOption(
// 		"View logs",
// 		"Open the logging report and inspect possible errors",
// 	),
// 	NewOption(
// 		"Quit",
// 		"Quit the application",
// 	),
// }

// // builds the header
// func (m model) headerView() string {
// 	// do more styling with lipgloss here...
// 	return "\n --------------- Claim Extraction Pipeline --------------- \n"
// } // headerView

// // build landing view
// func (m model) landingView() string {
// 	var s string
// 	landing := m.landing
// 	for _, stat := range landing.stats {
// 		desc := stat.description
// 		numFiles := stat.numFiles
// 		s += fmt.Sprintf("%s:\t%d\n", desc, numFiles)
// 	} // for
// 	s += landing.options.View()
// 	return s
// } // landingView

// // load landing page
// func initialModel() model {
// 	return model{
// 		currentView: landing,
// 		landing: LandingModel{
// 			stats: []Stat{
// 				NewStat(
// 					"Number of raw files",
// 					0,
// 				),
// 				NewStat(
// 					"Number of processed files",
// 					0,
// 				),
// 				NewStat(
// 					"Files pending extraction",
// 					0,
// 				),
// 			},
// 			options: NewOptionList(options),
// 		},
// 	}
// } // initialModel

// // GOAL: display landing page and make interactive
// // - banner
// // - stats
// //   - num of raw files
// //   - num of processed files
// //   - num of files pending extraction
// //
// // - options
// //   - run pipeline on new files
// //   - view processed files
// //   - quit

// func (m model) Init() tea.Cmd {
// 	return nil
// } // Init

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

// 	var cmd tea.Cmd
// 	var cmds []tea.Cmd

// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		k := msg.String()
// 		// user has quit
// 		if k == "ctrl+c" || k == "q" || k == "esc" {
// 			return m, tea.Quit
// 		} else if k == "enter" {
// 		} // if
// 	case tea.WindowSizeMsg:
// 		h, v := docStyle.GetFrameSize()
// 		m.landing.options.SetSize(msg.Width-h, msg.Height-v)
// 	} // switch

// 	m.landing.options, cmd = m.landing.options.Update(msg)
// 	cmds = append(cmds, cmd)
// 	return m, tea.Batch(cmds...)
// } // Update

// func (m model) View() string {
// 	// var b strings.Builder
// 	// var s string

// 	// return b.String()
// 	// return fmt.Sprintf("%s\n%s", m.headerView(), m.landingView())
// 	return lipgloss.Place(
// 		m.width,
// 		m.height,
// 		lipgloss.Center,
// 		lipgloss.Center,
// 		lipgloss.JoinVertical(
// 			lipgloss.Center,
// 			m.headerView(),
// 			m.landing.options.View(),
// 		),
// 	)
// } // View
