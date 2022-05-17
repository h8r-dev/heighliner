package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	hlnSectionStart = "# Added by hln"
	hlnSectionEnd   = "# End of section"
)

var hlnHostsSection = []string{
	"argocd",
	"nocalhost",
	"grafana",
	"prometheus",
	"loki",
}

type hostsOptions struct {
	genericclioptions.IOStreams
}

func newHostsCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := &hostsOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "hosts [appName]",
		Short: "Update hosts file",
		Args:  cobra.ExactArgs(1),
		RunE:  o.Run,
	}
	return cmd
}

func (o *hostsOptions) Run(cmd *cobra.Command, args []string) error {
	defaultHosts := filepath.Join("/etc", "hosts")
	appName := args[0]
	if _, err := getAppStatus(appName); err != nil {
		return fmt.Errorf("target app not found: %w", err)
	}
	hlnHostsSection = append(hlnHostsSection, appName)
	defaultTimeout := 90 * time.Second
	b, err := os.ReadFile(defaultHosts)
	if err != nil {
		return err
	}
	// Get ingress ip
	fmt.Fprintf(o.Out, "waiting for ingress pod to be ready...\n")
	if err := waitForIGController(defaultIngressNS, defaultIngressLabel, defaultTimeout); err != nil {
		return err
	}
	igip, err := getIngressIP(defaultIngressNS, defaultIngressSVC)
	if err != nil {
		return err
	}
	// Modify hosts
	lines := strings.Split(string(b), "\n")
	start, end, ok := findHlnSection(lines)
	newLines := getAppendHlnSection(igip)
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

func getAppendHlnSection(ingressIP string) []string {
	defaultSuffix := "h8r.site"
	lines := []string{}
	for _, prefix := range hlnHostsSection {
		line := ingressIP + " " + prefix + "." + defaultSuffix
		lines = append(lines, line)
	}
	return lines
}

// func validateHosts(ingressIP string, hlnLines []string) bool {

// 	return true
// }

// func updateHosts(ingressIP string, lines []string) bool {

// }
