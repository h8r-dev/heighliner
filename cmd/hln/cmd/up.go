package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

const upDesc = `
This command run a stack.

You should use '-s' or '--stack' to specify the stack. Use 'list stacks' subcommand 
to check all available stacks. Alternatively, you can use '--dir' flag 
to specify a local directory as your stack source. If you don't specify both '-s' 
and '--dir' flag, it will use current working directory by default:

    $ hln up -s gin-next

or

    $ hln up --dir /path/to/your/stack

To set values in a stack, use '-s' or '--stack' flag to specify a stack, use 
the '--set' flag and pass configuration from the command line:

    $ hln up -s gin-next --set foo=bar

You can specify the '--set' flag multiple times. The priority will be given to the
last (right-most) set specified. For example, if both 'bar' and 'newbar' values are
set for a key called 'foo', the 'newbar' value would take precedence:

    $ hln up -s gin-next --set foo=bar --set foo=newbar

Simply set '-i' or '--interactive' flag and it will run the stack interactively. You can 
fill your input values according to the prompts:

    $ hln up -s gin-next -i

`

// upOptions controls the behavior of up command.
type upOptions struct {
	Stack string
	Dir   string

	Values []string

	Interactive bool
	NoCache     bool

	genericclioptions.IOStreams
}

func (o *upOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.Stack, "stack", "s", "", "Name of your stack")
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVarP(&o.Interactive, "interactive", "i", false, "If this flag is set, heighliner will promt dialog when necessary.")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *upOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Stack != "" && o.Dir != "" {
		return errors.New("please do not specify both stack and dir")
	}
	for _, v := range o.Values {
		if !strings.Contains(v, "=") {
			return errors.New("value format should be '--set key=value'")
		}
	}
	return nil
}

func (o *upOptions) Run() error {
	// Save the pwd info brcause the program will chdir later.
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	// -----------------------------
	// 		Prepare stack
	// -----------------------------
	// Use local dir
	if o.Dir != "" {
		var err error
		o.Dir, err = homedir.Expand(o.Dir)
		if err != nil {
			return err
		}
	}
	// Use officaial stack
	if o.Stack != "" {
		stk, err := stack.New(o.Stack)
		if err != nil {
			return err
		}
		if err := stk.Update(); err != nil {
			return err
		}
		o.Dir = stk.Path
	}
	// -----------------------------
	//     	Set input values
	// -----------------------------
	// Handle --set flags
	for _, val := range o.Values {
		envvar := strings.Split(val, "=")
		envvar[1], err = homedir.Expand(envvar[1])
		if err != nil {
			return err
		}
		err := os.Setenv(envvar[0], envvar[1])
		if err != nil {
			return err
		}
	}
	// Handle interactive
	if o.Interactive {
		sch := schema.New()
		err = sch.AutomaticEnv(o.Interactive)
		if err != nil {
			return err
		}
	}
	// -----------------------------
	// 	Port-forward buildkit
	// -----------------------------
	// Forwarding port to buildkit
	readyCh := make(chan struct{})
	stopCh := make(chan struct{}, 1)
	errChan := make(chan error)
	port, err := util.GetAvailablePort()
	if err != nil {
		return err
	}

	go func() {
		errChan <- forwardPortToBuildKit(fmt.Sprintf("%d:%d", port, 1234), readyCh, stopCh)
	}()

	select {
	case <-readyCh:
		fmt.Println("PortForward to buildkit is ready")
	case err = <-errChan:
		fmt.Printf("PortForward to buildkit is terminated unexpectedly: %v\n", err)
		return err
	}
	_ = os.Setenv("BUILDKIT_HOST", fmt.Sprintf("tcp://127.0.0.1:%d", port))
	// -----------------------------
	// 	Execute dagger action
	// -----------------------------
	cli, err := dagger.NewClient(
		viper.GetString("log-format"),
		viper.GetString("log-level"),
		o.IOStreams,
	)
	if err != nil {
		return err
	}
	err = cli.Do(&dagger.ActionOptions{
		Name:    "up",
		Dir:     o.Dir,
		Plan:    "./plans",
		NoCache: o.NoCache,
	})
	if err != nil {
		return err
	}
	// -----------------------------
	// 	Handle the output
	// -----------------------------
	// Print the output.
	b, err := os.ReadFile(stackOutput)
	if err != nil {
		return err
	}
	fmt.Fprintf(o.IOStreams.Out, "%s\n", b)
	// Save the output info.
	err = copy.Copy(stackOutput, filepath.Join(pwd, appInfo))
	if err != nil {
		return err
	}

	return nil
}

func newUpCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &upOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Spin up your application",
		Long:  upDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(cmd, args); err != nil {
				return err
			}
			return o.Run()
		},
	}
	o.BindFlags(cmd.Flags())

	return cmd
}

func forwardPortToBuildKit(portStr string, readyCh, stopCh chan struct{}) error {
	fact := k8sutil.NewFactory(k8sutil.GetKubeConfigPath())
	client, err := fact.KubernetesClientSet()
	if err != nil {
		return err
	}

	// Find pod name of buildkit
	deploy, err := client.AppsV1().Deployments(heighlinerNs).Get(context.TODO(), buildKitName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	podList, err := client.CoreV1().Pods(heighlinerNs).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(deploy.Spec.Selector.MatchLabels).AsSelector().String()})
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return errors.New("no pod found for buildkit")
	}
	podName := podList.Items[0].Name // One pod only in this case

	restConfig, err := fact.ToRESTConfig()
	if err != nil {
		return err
	}

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(heighlinerNs).
		Name(podName).
		SubResource("portforward")

	iostream := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	transport, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())
	fw, err := portforward.NewOnAddresses(dialer, []string{"127.0.0.1"}, []string{portStr}, stopCh, readyCh, iostream.Out, iostream.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}
