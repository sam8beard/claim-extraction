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

const url = "https://charm.sh"

type model struct {
	status int
	err    error
} // model

type statusMsg int
type errMsg struct{ err error }

// for msgs that contain errors its usually handy to also implement
// the error interface on the message
func (e errMsg) Error() string { return e.err.Error() }

// this acts as a custom Cmd
func checkServer() tea.Msg {
	// create an http client and make get request
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return errMsg{err}
	} // if

	// we received a response from the server
	// return the http status code as a msg
	return statusMsg(res.StatusCode)

} // checkServer

func (m model) Init() tea.Cmd {
	// return the Cmd we made earlier
	// NOTE: the function is not called; the bubble tea runtime
	// 		 will do that when the time is right
	return checkServer
	// in bubbletea, Cmds return msgs
	// checkServer IS a Cmd that returns a msg
	// this Cmd's functionality is specified in the checkServer() function

} // Init

// func NewPrompt(prompt string) Prompt {
// 	return Prompt{prompt: prompt}
// }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		// the server returned a statusMsg
		// convert the statusMsg into an int
		// save the new int into the int var of the m struct
		// return the model and quit
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		// there was an error
		// save the errMsg into the err value of the m struct
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		// we received a key press
		// if the key press is ctrl+c, return the model and quit the program
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} // if
	} // switch

	// if we happen to get any other messages, just return the model
	// the model is the current state of our application, and we don't
	// want to change it in this case
	return m, nil
} // Update

// look at our current model and build
// the output string (the view of our application) accordingly
func (m model) View() string {
	// if theres an error, print it out and dont do anything
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	} // if

	// tell the user we're doing something
	s := fmt.Sprintf("Checking %s ...\n", url)

	// when the server responds with a status, add it to the current line
	if m.status > 0 {
		s += fmt.Sprintf("%d %s", m.status, http.StatusText(m.status))
	} // if

	// send off whatever we came up with above for rendering
	return "\n" + s + "\n\n"
} // View

func main() {
	_, err := tea.NewProgram(model{}).Run()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} // if

} // main
