package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/h8r-dev/heighliner/internal/k8sfactory"
)

const (
	hlnSectionStart   = "# Added by hln"
	hlnSectionEnd     = "# End of section"
	defaultIngressNS  = "ingress-nginx"
	defaultIngressSVC = "ingress-nginx-controller"
	defaultIngressIP  = "127.0.0.1" // This IP is for local kind and minikube
)

var hlnHostsSection = []string{
	"argocd",
	"nocalhost",
	"grafana",
	"prometheus",
	"loki",
}

type domainMappingOptions struct {
	IP     string
	Domain string

	genericclioptions.IOStreams
}

func (o *domainMappingOptions) BindFlags(f *pflag.FlagSet) {
	f.StringVar(&o.IP, "ip", "", "IP address")
	f.StringVar(&o.Domain, "domain", "", "Your domain name")
}

func newDomainMappingCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &domainMappingOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "domain-mapping [appName]",
		Short: "Set domain mapping",
		Args:  cobra.ExactArgs(1),
		RunE:  o.Run,
	}
	o.BindFlags(cmd.Flags())
	return cmd
}

func (o *domainMappingOptions) Run(cmd *cobra.Command, args []string) error {
	defaultHosts := filepath.Join("/etc", "hosts")
	appName := args[0]
	// This section will test if the app exists.
	// if _, err := getAppStatus(appName); err != nil {
	// 	return fmt.Errorf("target app not found: %w", err)
	// }
	hlnHostsSection = append(hlnHostsSection, appName)
	b, err := os.ReadFile(defaultHosts)
	if err != nil {
		return err
	}
	// Get ingress ip
	ip := ""
	if o.IP != "" {
		ip = o.IP
	} else {
		igip, err := getIngressIP(defaultIngressNS, defaultIngressSVC)
		if err != nil {
			return err
		}
		ip = igip
	}
	domain := "h8r.site"
	if o.Domain != "" {
		domain = o.Domain
	}
	// Modify hosts
	lines := strings.Split(string(b), "\n")
	start, end, ok := findHlnSection(lines)
	newLines := getAppendHlnSection(ip, domain)
	if !ok {
		if err := copy.Copy(defaultHosts, defaultHosts+".bak"); err != nil {
			return fmt.Errorf("failed to backup hosts file: %w", err)
		}
		lines = append(lines, hlnSectionStart)
		lines = append(lines, newLines...)
		lines = append(lines, hlnSectionEnd)
		data := []byte(strings.Join(lines, "\n"))
		if err := os.WriteFile(defaultHosts, data, 0644); err != nil {
			return fmt.Errorf("failed to write hosts file: %w", err)
		}
		return nil
	}
	// Update hosts
	headLines := lines[:start+1]
	tailLines := lines[end:]
	contents := []string{}
	contents = append(contents, headLines...)
	contents = append(contents, newLines...)
	contents = append(contents, tailLines...)
	data := []byte(strings.Join(contents, "\n"))
	if err := os.WriteFile(defaultHosts, data, 0644); err != nil {
		return fmt.Errorf("failed to write hosts file: %w", err)
	}
	return nil
}

func findHlnSection(lines []string) (start, end int, ok bool) {
	start, end = -1, -1
	for i, line := range lines {
		if line == hlnSectionStart {
			start = i
			break
		}
	}
	if start < 0 {
		return
	}
	for i, line := range lines[start:] {
		if line == hlnSectionEnd {
			end = start + i
			break
		}
	}
	if start < end {
		ok = true
	}
	return
}

func getAppendHlnSection(ip, domain string) []string {
	lines := []string{}
	for _, prefix := range hlnHostsSection {
		line := ip + " " + prefix + "." + domain
		lines = append(lines, line)
	}
	return lines
}

// func validateHosts(ingressIP string, hlnLines []string) bool {

// 	return true
// }

// func updateHosts(ingressIP string, lines []string) bool {

// }

func getIngressIP(namespace, svcName string) (string, error) {
	cs, err := k8sfactory.GetDefaultClientSet()
	if err != nil {
		return "", err
	}
	ctx := context.TODO()
	igsvc, err := cs.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	igs := igsvc.Status.LoadBalancer.Ingress
	var igip string
	if len(igs) > 0 {
		igip = igs[0].IP
	} else {
		igip = defaultIngressIP
	}
	return igip, err
}
