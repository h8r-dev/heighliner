package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/h8r-dev/heighliner/pkg/state/app"
)

// LogsOptions controls the behavior of logs command.
type LogsOptions struct {
	Choice int

	// PodLogOptions
	Follow bool

	Kubecli *kubernetes.Clientset
}

func newLogsCmd() *cobra.Command {
	o := &LogsOptions{}

	cmd := &cobra.Command{
		Use:   "logs [appName]",
		Args:  cobra.ExactArgs(1),
		Short: "Print the logs for an app",
		RunE:  o.getPodLogs,
	}
	o.addFlags(cmd)
	return cmd
}

func (o *LogsOptions) addFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", o.Follow, "Specify if the logs should be streamed.")
}

func getServiceNames(services []app.Service) []string {
	var names []string
	for _, s := range services {
		names = append(names, s.Name)
	}
	return names
}

func (o *LogsOptions) getPodLogs(cmd *cobra.Command, args []string) error {

	k8sClient, err := getDefaultClientSet()
	if err != nil {
		return err

	}
	o.Kubecli = k8sClient

	st, err := getStateInSpecificBackend()
	if err != nil {
		return err
	}
	appInfo, err := st.LoadOutput(args[0])
	if err != nil {
		return err
	}

	names := getServiceNames(appInfo.Services)

	// ask user to select one of the services to get logs from
	p := tea.NewProgram(initialModel(names, &o.Choice))
	if err := p.Start(); err != nil {
		return err
	}

	namespace := fmt.Sprintf("%s-deploy-production", appInfo.ApplicationRef.Name)
	svc, err := o.Kubecli.CoreV1().Services(namespace).Get(context.TODO(), names[o.Choice], metav1.GetOptions{})
	if err != nil {
		return err
	}

	podlist, err := o.Kubecli.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&metav1.LabelSelector{
			MatchLabels: svc.Spec.Selector,
		}),
	})
	if err != nil {
		return err
	}
	if len(podlist.Items) == 0 {
		return fmt.Errorf("no pods found for service %s", names[o.Choice])
	}

	podNames := []string{}
	for _, po := range podlist.Items {
		podNames = append(podNames, po.Name)
	}

	var podChoice int
	// ask user to select one of the services to get logs from
	program := tea.NewProgram(initialModel(podNames, &podChoice))
	if err := program.Start(); err != nil {
		return err
	}

	request := o.Kubecli.CoreV1().Pods(namespace).GetLogs(podNames[podChoice], &corev1.PodLogOptions{
		Follow: o.Follow,
	})

	return DefaultConsumeRequest(request, os.Stdout)
}

// DefaultConsumeRequest reads the data from request and writes into
// the out writer. It buffers data from requests until the newline or io.EOF
// occurs in the data, so it doesn't interleave logs sub-line
// when running concurrently.
func DefaultConsumeRequest(request rest.ResponseWrapper, out io.Writer) error {
	readCloser, err := request.Stream(context.TODO())
	if err != nil {
		return err
	}
	defer func() {
		if err := readCloser.Close(); err != nil {
			panic(err)
		}
	}()

	r := bufio.NewReader(readCloser)
	for {
		bytes, err := r.ReadBytes('\n')
		if _, err := out.Write(bytes); err != nil {
			return err
		}

		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}

type model struct {
	choices   []string // items on the to-do list
	cursor    int      // which to-do list item our cursor is pointing at
	choiceRef *int
}

func initialModel(choices []string, choiceRef *int) model {
	return model{
		choices:   choices,
		choiceRef: choiceRef,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// nolint
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			*m.choiceRef = m.cursor
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select a service to get logs from\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
