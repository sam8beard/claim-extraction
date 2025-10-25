package main

import (
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/bubbles/textinput"
	// "github.com/charmbracelet/lipgloss"
	"fmt"
)

const charmUrl = "https://charm.sh"
const googleUrl = "https://google.com"

type model struct {
	status     int
	err        error
	receivedAt time.Time
	state      string
	url        string
	choice     string
} // model

type loadingMsg string
type statusMsg int
type errMsg struct{ err error }

// type stateMsg string

// type readyMsg tea.Msg

// for msgs that contain errors its usually handy to also implement
// the error interface on the message
func (e errMsg) Error() string { return e.err.Error() }

func initialModel() model {
	return model{state: "landing"}
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
	return nil
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
	load := createLoader()

	if msg != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// change state to loading
			m.state = "loading"
			switch msg.String() {
			case "enter":
				m.choice = "Charm"
				return m, checkCharmServer
			case " ":
				m.choice = "Google"
				return m, checkGoogleServer
			case "q":
				m.state = "quit"
				return m, tea.Quit

			case "b": // testing return on error
				return m, checkServer("https://harm.sh")
			} // switch

		case statusMsg:
			load()
			m.state = "ready"
			m.status = int(msg)
			m.receivedAt = time.Now()
			return m, nil
		case errMsg:
			load()
			m.state = "ready"
			m.err = msg
			return m, tea.Quit
		default:
			return m, nil
		} // switch

	} // if

	return m, nil
	// switch msg := msg.(type) {
	// case statusMsg:
	// 	// the server returned a statusMsg
	// 	// convert the statusMsg into an int
	// 	// save the new int into the int var of the m struct
	// 	// return the model and quit
	// 	m.state = "get"
	// 	m.status = int(msg)
	// 	m.receivedAt = time.Now()
	// 	return m, initialModel
	// case errMsg:
	// 	// there was an error
	// 	// save the errMsg into the err value of the m struct
	// 	m.err = msg
	// 	return m, tea.Quit
	// case tea.KeyMsg:
	// 	// we received a key press
	// 	// if the key press is ctrl+c, return the model and quit the program

	// 	if msg.Type == tea.KeyCtrlC {
	// 		return m, tea.Quit
	// 	} // if

	// 	if msg.Type == tea.KeyEnter {
	// 		return m, checkCharmServer
	// 	}

	// 	if msg.Type == tea.KeySpace {
	// 		return m, checkGoogleServer
	// 	}
	// } // switch

	// if we happen to get any other messages, just return the model
	// the model is the current state of our application, and we don't
	// want to change it in this case

} // Update

// look at our current model and build
// the output string (the view of our application) accordingly
func (m model) View() string {
	// if theres an error, print it out and dont do anything
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	} // if

	// s := "\nPress enter for Charm and space for Google\n"
	var s string
	switch m.state {
	case "landing":
		return "\nPress enter for Charm and space for Google\n"
	case "loading":
		s = fmt.Sprintf("\nLoading %s...\n", m.choice)
		return s
	case "ready":
		s = fmt.Sprintf("%d %s \nReceived at: %s", m.status, http.StatusText(m.status), m.receivedAt)
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
