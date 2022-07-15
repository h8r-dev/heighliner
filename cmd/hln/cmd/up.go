package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/kubectl/pkg/cmd/config"
	"k8s.io/kubectl/pkg/scheme"

	"github.com/h8r-dev/heighliner/pkg/dagger"
	"github.com/h8r-dev/heighliner/pkg/logger"
	"github.com/h8r-dev/heighliner/pkg/schema"
	"github.com/h8r-dev/heighliner/pkg/stack"
	"github.com/h8r-dev/heighliner/pkg/state"
	"github.com/h8r-dev/heighliner/pkg/util"
	"github.com/h8r-dev/heighliner/pkg/util/cueutil"
	"github.com/h8r-dev/heighliner/pkg/util/k8sutil"
)

const upDesc = `
This command run a stack.

You should use '-s' or '--stack' to specify the stack. Use 'list stacks' subcommand 
to check all available stacks. Alternatively, you can use '--dir' flag 
to specify a local directory as your stack source. If you don't specify both '-s' 
and '--dir' flag, it will use current working directory by default:

    $ hln up [appName] -s gin-next

or

    $ hln up [appName] --dir /path/to/your/stack

To set values in a stack, use '-s' or '--stack' flag to specify a stack, use 
the '--set' flag and pass configuration from the command line:

    $ hln up [appName] -s gin-next --set foo=bar

You can specify the '--set' flag multiple times. The priority will be given to the
last (right-most) set specified. For example, if both 'bar' and 'newbar' values are
set for a key called 'foo', the 'newbar' value would take precedence:

    $ hln up [appName] -s gin-next --set foo=bar --set foo=newbar

Simply set '-i' or '--interactive' flag and it will run the stack interactively. You can 
fill your input values according to the prompts:

    $ hln up [appName] -s gin-next -i

`

// upOptions controls the behavior of up command.
type upOptions struct {
	Stack   string
	Version string
	Dir     string

	File   string
	Values []string

	Interactive bool
	NoCache     bool

	genericclioptions.IOStreams
}

func (o *upOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVarP(&o.Stack, "stack", "s", "", "Name of your stack")
	f.StringVarP(&o.File, "file", "f", "", "Path to your input file")
	f.StringVar(&o.Dir, "dir", "", "Path to your local stack")
	f.StringArrayVar(&o.Values, "set", []string{}, "The input values of your project")
	f.BoolVarP(&o.Interactive, "interactive", "i", false, "If this flag is set, heighliner will prompt dialog when necessary.")
	f.BoolVar(&o.NoCache, "no-cache", false, "Disable caching")
}

func (o *upOptions) Validate(cmd *cobra.Command, args []string) error {
	if o.Stack != "" && o.Dir != "" {
		return errors.New("can't use both stack and dir")
	}
	for _, v := range o.Values {
		if !strings.Contains(v, "=") {
			return errors.New("format of values should be '--set key=value'")
		}
	}
	return nil
}

func (o *upOptions) Complete() error {
	if strings.Contains(o.Stack, "@") {
		args := strings.Split(o.Stack, "@")
		if len(args) < 2 {
			return errors.New("invalid stack format, should be name@version")
		}
		o.Stack = args[0]
		o.Version = args[1]
	}
	return nil
}

func (o *upOptions) Run() error {
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
		stk, err := stack.New(o.Stack, o.Version)
		if err != nil {
			return err
		}
		if err := stk.Update(); err != nil {
			return err
		}
		o.Dir = stk.Path
	}
	// -----------------------------
	//     	Convert input file
	// -----------------------------
	if o.File != "" {
		if err := cueutil.ConvertYamlToCue(o.File, filepath.Join(o.Dir, "plans", "input.cue")); err != nil {
			return fmt.Errorf("failed to convert input file: %w", err)
		}
	}
	// -----------------------------
	//     	Set input values
	// -----------------------------
	// Handle --set flags
	if err := o.setEnv(); err != nil {
		return err
	}
	// -----------------------------
	// 	Port-forward buildkit
	// -----------------------------
	// Forwarding port to buildkit
	if err := runForward(o.IOStreams); err != nil {
		return err
	}
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

	fmt.Fprintf(o.Out, "\n%s\n", color.GreenString("🎉 Congrats! Application is ready!"))
	return nil
}

func (o upOptions) setEnv() error {
	for _, val := range o.Values {
		envs := strings.Split(val, "=")
		if len(envs) != 2 {
			return errors.New("value format should be '--set key=value'")
		}
		key, val := envs[0], envs[1]
		val, err := homedir.Expand(val)
		if err != nil {
			return err
		}
		val, err = filepath.Abs(val)
		if err != nil {
			return err
		}
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	// Handle interactive
	if o.Interactive {
		sch := schema.New(o.Dir)
		if err := sch.AutomaticEnv(o.Interactive); err != nil {
			return err
		}
	}
	return nil
}

func newUpCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &upOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "up [appName]",
		Short: "Spin up your application",
		Long:  upDesc,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Validate(cmd, args); err != nil {
				return err
			}
			return o.Complete()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}
	o.BindFlags(cmd.Flags())

	return cmd
}

func runForward(streams genericclioptions.IOStreams) error {
	lg := logger.New(streams)
	readyCh := make(chan struct{})
	stopCh := make(chan struct{}, 1)
	errChan := make(chan error)
	port, err := util.GetAvailablePort()
	if err != nil {
		return err
	}

	go func() {
		errChan <- forwardPortToBuildKit(streams, fmt.Sprintf("%d:%d", port, 1234), readyCh, stopCh)
	}()

	select {
	case <-readyCh:
		lg.Info("port-forward to buildkit is ready")
	case err = <-errChan:
		return fmt.Errorf("port-forward to buildkit is terminated unexpectedly: %w", err)
	}

	if err := os.Setenv("BUILDKIT_HOST", fmt.Sprintf("tcp://127.0.0.1:%d", port)); err != nil {
		return err
	}

	lg.Info(fmt.Sprintf("flattening kubeconfig: %s\n", k8sutil.GetKubeConfigPath()))
	if err := flattenKubeconfig(); err != nil {
		return fmt.Errorf("failed to flatten kubeconfig: %w", err)
	}
	return nil
}

func forwardPortToBuildKit(streams genericclioptions.IOStreams, portStr string, readyCh, stopCh chan struct{}) error {
	fact := k8sutil.NewFactory(k8sutil.GetKubeConfigPath())
	client, err := fact.KubernetesClientSet()
	if err != nil {
		return err
	}

	// Find pod name of buildkit
	deploy, err := client.AppsV1().Deployments(state.HeighlinerNs).Get(context.TODO(), buildKitName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	podList, err := client.CoreV1().Pods(state.HeighlinerNs).List(context.TODO(), metav1.ListOptions{
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
		Namespace(state.HeighlinerNs).
		Name(podName).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(restConfig)
	if err != nil {
		return err
	}
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())
	fw, err := portforward.NewOnAddresses(dialer, []string{"127.0.0.1"}, []string{portStr}, stopCh, readyCh, streams.Out, streams.ErrOut)
	if err != nil {
		return err
	}
	return fw.ForwardPorts()
}

func flattenKubeconfig() error {
	kubeconfig := k8sutil.GetKubeConfigPath()

	b := make([]byte, 0)
	buff := bytes.NewBuffer(b)
	po := clientcmd.NewDefaultPathOptions() // po.LoadingRules.ExplicitPath = kubeconfig
	vo := config.ViewOptions{
		ConfigAccess: po,
		Flatten:      true,
		Merge:        1,
		PrintFlags:   genericclioptions.NewPrintFlags("").WithTypeSetter(scheme.Scheme).WithDefaultOutput("yaml"),
		Minify:       true,
		IOStreams:    genericclioptions.IOStreams{In: os.Stdin, Out: buff, ErrOut: os.Stderr},
	}
	printer, err := vo.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	vo.PrintObject = printer.PrintObj
	if err = vo.Run(); err != nil {
		return err
	}

	bys, err := ioutil.ReadFile(kubeconfig)
	if err != nil {
		return err
	}

	// Bak kubeconfig
	if err = ioutil.WriteFile(kubeconfig+".hln.bak", bys, 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(kubeconfig, buff.Bytes(), 0644)
}
