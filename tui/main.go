package main 

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"fmt"
	"os"
	"log"
)
type Styles struct { 
	BorderColor lipgloss.Color
	InputField lipgloss.Style

}

func DefaultStyles() *Styles { 
	s := new(Styles)
	s.BorderColor = lipgloss.Color("42")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct { 
	prompts []Prompt // search topics
	answerField textinput.Model // input field
	width int 
	height int
	index int // what prompt we are currently on
	styles *Styles
	
} // model

type Prompt struct { 
	prompt string
	answer string
}

func NewPrompt(prompt string) Prompt { 
	return Prompt{prompt: prompt}
}

func New(prompts []Prompt) *model { 
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = ""
	answerField.Focus()
	return &model{
		prompts: prompts, 
		answerField: answerField, 
		styles: styles,
	}

} // New


func (m model) Init() tea.Cmd { 
	// just return nil, this means "not accepting I/O right now"
	return nil
} // Init

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { 
	var cmd tea.Cmd
	current := &m.prompts[m.index]
	switch msg  := msg.(type) {

		case tea.WindowSizeMsg: 
			m.width = msg.Width 
			m.height = msg.Height

		// Is it a key press? 
		case tea.KeyMsg:

			// what key was pressed?
			switch msg.String() { 
			
			// keys that exit a program
			case "ctrl+c", "q": 
				return m, tea.Quit
			case "enter": 
				current.answer = m.answerField.Value()
				m.answerField.SetValue("")
				log.Printf("prompt: %s, answer: %s", current.prompt, current.answer)
				m.Next()
 				return m, nil 
		} // switch
	} // switch
	
	m.answerField, cmd = m.answerField.Update(msg)
	 // Return the model to the BT runtime for processing 
	 // NOTE that we are not returning a command 
	 return m, cmd
} // Update

func (m model) View() string { 

	if m.width == 0 { 
		return "loading..."
	} // if 
	
	return lipgloss.Place(
		m.width, 
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center, 
			m.prompts[m.index].prompt, 
			m.styles.InputField.Render(m.answerField.View()), 
		), 
	)
	

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

func (m *model) Next() { 
	if m.index < len(m.prompts)-1 { 
		m.index++
	} else { 
		m.index = 0
	}
}
func main() { 
	// set up error logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil { 
			log.Fatalf("err: %w", err)
	} // if 
	defer f.Close()
	
	prompts := []Prompt{
		NewPrompt("What topic would you like to search?"), 
		NewPrompt("How many files would you like to analyze?"),
		NewPrompt("What would you like the confidence threshold to be?"),
	}
	m := New(prompts)
	
	
	
	// make program
	p := tea.NewProgram(m, tea.WithAltScreen())

	// start program
	if _, err := p.Run(); err != nil { 
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	} // if 
} // main