package main 

import (
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/lipgloss"
	"fmt"
	"os"
	"log"
)
// type Styles struct { 
// 	BorderColor lipgloss.Color
// 	InputField lipgloss.Style

// }

// func DefaultStyles() *Styles { 
// 	s := new(Styles)
// 	s.BorderColor = lipgloss.Color("42")
// 	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
// 	return s
// }

type model struct { 
	options []string
	cursor int 
	selected string
} // model

func initialModel() model { 
	return model{ 
		options: []string{
			"View processed files",
			"Run pipeline",
			"View logs",
		}, 
		cursor: 0,
		selected: "",
	}
} // initialModel

func (m model) Init() tea.Cmd { 
	// just return nil, this means "not accepting I/O right now"
	return nil
} // Init


// func NewPrompt(prompt string) Prompt { 
// 	return Prompt{prompt: prompt}
// }

// func New(prompts []Prompt) *model { 
// 	styles := DefaultStyles()
// 	answerField := textinput.New()
// 	answerField.Placeholder = ""
// 	answerField.Focus()
// 	return &model{
// 		prompts: prompts, 
// 		answerField: answerField, 
// 		styles: styles,
// 	}

// } // New




func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { 
	// var cmd tea.Cmd
	// current := &m.prompts[m.index]
	switch msg  := msg.(type) {

		// case tea.WindowSizeMsg: 
		// 	m.width = msg.Width 
		// 	m.height = msg.Height

		// Is it a key press? 
		case tea.KeyMsg:
			// what key was pressed?
			switch msg.String() { 
			
			// keys that exit a program
			case "ctrl+c", "q": 
				return m, tea.Quit
			case "up":
				// if cursor is not already at first option
				if m.cursor > 0 { 
					m.cursor--
				} // if
			case "down": 
				// if cursor is not already at last option
				if m.cursor < len(m.options)-1 {
					m.cursor++
				} // if 
			case "enter": 
				m.selected = m.options[m.cursor]
				return m, nil
		} // switch
	} // switch
	
	 // Return the model to the BT runtime for processing 
	 // NOTE that we are not returning a command 
	 return m, nil
} // Update

func (m model) View() string { 

	// if m.width == 0 { 
	// 	return "loading..."
	// } // if 
	s:= "--------- Claim Extraction Pipeline ---------\n\n"
	for i, option := range m.options { 
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">"
		} // if 
		s += fmt.Sprintf("%s %s\n", cursor, option)
	} // for 
	if m.selected != "" { 
		s += fmt.Sprintf("\nYou chose: %s\n", m.selected)
	} // if 
	s += "\nPress q to quit.\n"
	return s
	
	

	// // header
	// s := "What would you like to do today?\n\n"

	// // iterate over choices 
	// for i, choice := range m.choices { 

	// 	// is the cursor pointing at this choice? 
	// 	cursor := " " // no cursor
	// 	if m.cursor == i { 
	// 		cursor = ">" // cursor
	// 	} // if 

	// 	// is this choice selected?
	// 	checked := " " // not selected
	// 	if _, ok := m.selected[i]; ok { 
	// 		checked = "x" // selected
	// 	} // if 

	// 	// render row
	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	// } // for 

	// // footer 
	// s += "\nPress q to quit.\n"

	// // send UI for rendering 
	// return s
} // View

// func (m *model) Next() { 
// 	if m.index < len(m.prompts)-1 { 
// 		m.index++
// 	} else { 
// 		m.index = 0
// 	}
// }
func main() { 
	// set up error logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil { 
			log.Fatalf("err: %w", err)
	} // if 
	defer f.Close()

	
	// make program
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	// start program
	if _, err := p.Run(); err != nil { 
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	} // if 
} // main