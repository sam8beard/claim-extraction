package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	// "github.com/charmbracelet/lipgloss"
	"fmt"
)

type model struct {
	status         int             // response status
	err            error           // error received
	receivedAt     time.Time       // time response was received
	state          string          // state of model
	url            string          // url to visit
	textInput      textinput.Model // text input for url
	visits         map[string]int  // visits to each site
	acceptingInput bool            // indicates whether or not input is being accepted
} // model

type statusMsg int
type errMsg struct{ err error }

// for msgs that contain errors its usually handy to also implement
// the error interface on the message
func (e errMsg) Error() string { return e.err.Error() }

func buildTextInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Search for a website"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return ti
} // buildTextInput

func initialModel() model {
	return model{
		acceptingInput: true,
		textInput:      buildTextInput(),
		state:          "",
		visits: map[string]int{
			"Google.com": 0,
			"Charm.sh":   0,
			"Harm.sh":    0,
		},
	}
} // initialModel

func checkServer(url string) tea.Cmd {
	// load := createLoader()
	return func() tea.Msg {
		// load()
		c := &http.Client{Timeout: 10 * time.Second}

		res, err := c.Get(url)
		if err != nil {
			return errMsg{err}
		}
		return statusMsg(res.StatusCode)
	}
}

func (m model) Init() tea.Cmd {

	return textinput.Blink

} // Init

// simulate a loading state
func createLoader() func() {
	return func() {
		duration := time.Duration(1)
		time.Sleep(duration * time.Second)
	}
} // createLoader

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// create a load simulator
	load := createLoader()

	// create cmd for text input
	var cmd tea.Cmd

	// switch on message type
	switch msg := msg.(type) {
	// if the message is a key press
	case tea.KeyMsg:
		// switch on key type
		switch msg.Type {

		// user submits input
		case tea.KeyEnter:
			// disable input
			m.textInput.Blur()
			m.acceptingInput = false
			m.state = "submit"
			m.url = m.textInput.Value()
			m.textInput.SetValue("")
			return m, checkServer(m.url)

		// user quits
		case tea.KeyCtrlC:
			m.state = "quit"
			return m, tea.Quit

		// attempt to let user type
		default:
			// accept input if user has not just submitted something
			if m.acceptingInput {
				m.textInput.Focus()
				m.textInput, cmd = m.textInput.Update(msg)
			} // if
		} // switch

	// msg is valid response
	case statusMsg:
		// simulate load
		load()
		m.acceptingInput = true
		m.state = "received"
		m.status = int(msg)
		m.receivedAt = time.Now()
		return m, nil

	// msg is invaild response
	case errMsg:
		// simulate load
		load()
		m.acceptingInput = true
		m.state = "error"
		m.err = msg
		return m, nil

	// if nothing, wait for message
	default:
		return m, nil
	} // switch
	return m, cmd
} // Update

// STATES

// submit: shows loading screen
// - msg received: tea.KeyType == "enter"
// - m.state: change state to submit
// - cmd returned: checkServer(m.url)
// - this will be shown when the msg received is cmd returned is checkServer

// ** i think loading needs to happen here **
// received: shows response
// - msg received: statusMsg
// - m.state: change state to received
// - cmd returned: nil

// error: shows error
// - msg received: errMsg
// - m.state: change state to error
// - cmd returned: nil

// quit: shows quit message
// - msg received: tea.KeyMsg == "ctrl+c"
// - m.state: change state to quit
// - cmd returned: tea.Quit

// look at our current model and build
// the output string (the view of our application) accordingly
func (m model) View() string {
	var b strings.Builder
	var s string
	switch m.state {
	// user submits url
	case "submit":
		s = fmt.Sprintf("\nLoading %s...\n", m.url)
		b.WriteString(s)
	// print response received
	case "received":
		s = fmt.Sprintf("%d %s \nReceived at: %s\n", m.status, http.StatusText(m.status), m.receivedAt)
		b.WriteString(s)
	// print error
	case "error":
		s = fmt.Sprintf("\nWebsite not found: %v\n\n", m.err)
		b.WriteString(s)
	// user quit
	case "quit":
		return "\nQuitting application\n"
	} // switch

	// show prompt when we are accepting input
	if m.acceptingInput {
		s = ("Enter a url: " + m.textInput.View())
		b.WriteString("\n" + s + "\n\n")
	} // if

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
