package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var words []string = []string{}
var s = rand.NewSource(time.Now().UnixNano())
var r = rand.New(s)

func main() {
	f, err := os.Open("words.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// Add the word to a slice of strings
		words = append(words, scanner.Text())
	}

	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type tickMsg struct{}
type errMsg error

type model struct {
	textInput input.Model
	word      int
	err       error
}

func initialModel() model {
	inputModel := input.NewModel()
	inputModel.Placeholder = "type here"
	inputModel.Focus()

	return model{
		textInput: inputModel,
		word:      r.Intn(len(words)),
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	return input.Blink(m.textInput)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.textInput.Reset()
			m.word = r.Intn(len(words))
			return m, cmd
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = input.Update(msg, m.textInput)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n> %s\n\n\n\n%s",
		input.View(m.textInput),
		words[m.word],
		"(esc to quit)",
	) + "\n"
}
