package schema

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

	"github.com/h8r-dev/heighliner/pkg/util"
)

func startUI(pm Parameter) error {
	p := tea.NewProgram(initialModel(pm))
	m, err := p.StartReturningModel()
	if err != nil {
		return err
	}
	if m, ok := m.(model); ok && errors.Is(m.err, ErrCancelInput) {
		return m.err
	}
	return nil
}

// setval will be called when user presses enter.
func setVal(p Parameter, val string) error {
	switch {
	case val != "":
		if err := os.Setenv(p.Key, util.Abs(strings.TrimSpace(val))); err != nil {
			panic(err)
		}
	case !p.Required:
		return nil
	default:
		return errValueMissed
	}
	return nil
}

// ------
// Logic of Terminal UI
// ------

var (
	errValueMissed = errors.New("this value is required")
	// ErrCancelInput is a signal to break the interactive inputing process.
	ErrCancelInput = errors.New("cancel interactive inputing process")
)

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
	ti.CharLimit = 0
	ti.Width = 50

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
		case tea.KeyCtrlC:
			m.err = ErrCancelInput
			return m, tea.Quit
		case tea.KeyTab:
			m.textInput.SetValue(m.parameter.Default)
			return m, nil
		case tea.KeyEnter:
			if err := setVal(m.parameter, m.textInput.Value()); err != nil {
				m.err = err
				return m, nil
			}
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
	s := m.parameter.Description
	if m.parameter.Required {
		s += color.YellowString(" (required)")
	}
	s += fmt.Sprintf(
		": \n\n%s\n\n",
		m.textInput.View(),
	)
	if errors.Is(m.err, errValueMissed) {
		s += color.RedString("Warn: %s", m.err.Error())
	}
	return s
}
