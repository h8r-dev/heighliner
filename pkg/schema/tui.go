package schema

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

func startui(pm Parameter) error {
	p := tea.NewProgram(initialModel(pm))

	return p.Start()
}

// setval will be called when user presses enter.
func setVal(p Parameter, val string) {
	switch {
	case val != "":
		if err := os.Setenv(p.Key, val); err != nil {
			panic(err)
		}
	case p.Default != "":
		if err := os.Setenv(p.Key, val); err != nil {
			panic(err)
		}
	default:
		color.Red("\nvalue couldn't be empty")
		os.Exit(1)
	}
}

// ------
// Logic of Terminal UI
// ------

type errMsg error

type model struct {
	textInput textinput.Model
	parameter Parameter
	err       error
}

func initialModel(p Parameter) model {
	ti := textinput.New()
	ti.Placeholder = p.Default
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		parameter: p,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			setVal(m.parameter, m.textInput.Value())
			return m, tea.Quit
		default:
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Please input %s (%s):\n\n%s\n\n%s",
		m.parameter.Title,
		m.parameter.Description,
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}
