package main

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	// "github.com/charmbracelet/lipgloss"
	"fmt"
)

const charmUrl = "https://charm.sh"
const googleUrl = "https://google.com"

type model struct {
	status     int             // response status
	err        error           // error received
	receivedAt time.Time       // time response was received
	state      string          // state of model
	url        string          // url to visit
	textInput  textinput.Model // text input for url
	visits     map[string]int  // visits to each site
	choice     string          // site chosen by keypress
	warning    bool            // indicator for malformed url
} // model

type loadingMsg string
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
		textInput: buildTextInput(),
		state:     "landing",
		visits: map[string]int{
			"Google.com": 0,
			"Charm.sh":   0,
			"Harm.sh":    0,
		},
	}
} // initialModel

// this acts as a custom Cmd
func checkCharmServer() tea.Msg {
	// create an http client and make get request
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(charmUrl)
	if err != nil {
		return errMsg{err}
	} // if

	// we received a response from the server
	// return the http status code as a msg
	return statusMsg(res.StatusCode)

} // checkCharmServer

func checkServer(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)
		if err != nil {
			return errMsg{err}
		}
		return statusMsg(res.StatusCode)
	}
}

func checkGoogleServer() tea.Msg {
	// create an http client and make get request
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(googleUrl)
	if err != nil {
		return errMsg{err}
	} // if

	// we received a response from the server
	// return the http status code as a msg
	return statusMsg(res.StatusCode)
} // checkBadServer

func (m model) Init() tea.Cmd {
	// return the Cmd we made earlier
	// NOTE: the function is not called; the bubble tea runtime
	// 		 will do that when the time is right

	// in bubbletea, Cmds return msgs
	// checkCharmServer IS a Cmd that returns a msg
	// this Cmd's functionality is specified in the checkCharmServer() function
	return textinput.Blink
} // Init

//	func NewPrompt(prompt string) Prompt {
//		return Prompt{prompt: prompt}
//	}
//
//	func ready() tea.Msg {
//		return readyMsg("ready")
//	} // ready

// simulate a loading state
func createLoader() func() {
	return func() {
		duration := time.Duration(2)
		time.Sleep(duration * time.Second)
	}
} // createLoader

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// create a load simulator
	load := createLoader()

	if msg != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// change state to loading
			m.state = "loading"
			switch msg.String() {
			case "enter":
				m.choice = "Charm.sh"
				return m, checkCharmServer
			case " ":
				m.choice = "Google.com"
				return m, checkGoogleServer
			case "b": // testing return on error
				m.choice = "Harm.sh"
				return m, checkServer("https://harm.sh")
			case "q":
				m.state = "quit"
				return m, tea.Quit
			default:
				m.state = "quit"
				return m, tea.Quit
			} // switch

		case statusMsg:
			load()
			m.warning = false
			m.visits[m.choice]++
			m.state = "ready"
			m.status = int(msg)
			m.receivedAt = time.Now()
			return m, nil
		case errMsg:
			load()
			m.warning = true
			m.state = "landing"
			m.err = msg
			return m, nil
		default:
			return m, nil
		} // switch

	} // if

	return m, nil

} // Update

// look at our current model and build
// the output string (the view of our application) accordingly
func (m model) View() string {
	// if theres an error, print it out and dont do anything
	var s string
	if m.warning {
		s = fmt.Sprintf("\nWebsite not found: %v\n\n", m.err)
	} // if

	// s := "\nPress enter for Charm and space for Google\n"
	switch m.state {
	case "landing":
		s += "\nPress enter for Charm and space for Google\n"
	case "loading":
		s = fmt.Sprintf("\nLoading %s...\n", m.choice)
		return s
	case "ready":
		s = fmt.Sprintf("%d %s \nReceived at: %s\n", m.status, http.StatusText(m.status), m.receivedAt)
		s += fmt.Sprintf(
			"\nVisits: \n%s | %d\n%s | %d\n%s | %d\n",
			"Google.com",
			m.visits["Google.com"],
			"Charm.sh",
			m.visits["Charm.sh"],
			"Harm.sh",
			m.visits["Harm.sh"],
		)

	case "quit":
		return "\nQuitting application\n"
	} // switch

	// // tell the user we're doing something
	// s += fmt.Sprintf("Checking %s ...\n", charmUrl)

	// // when the server responds with a status, add it to the current line
	// if m.status > 0 {
	// 	s += fmt.Sprintf("%d %s \nReceived at: %s", m.status, http.StatusText(m.status), m.receivedAt)
	// } // if

	// send off whatever we came up with above for rendering
	return "\n" + s + "\n\n"
} // View

func main() {
	p := tea.NewProgram(initialModel())

	_, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} // if

} // main
