package terraform

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-exec/tfexec"
)

type ApplyOptions struct {
	Dir        string
	Suffix     string
	Namespace  string
	KubeConfig string
}

func NewApplyOptions(dir, suffix, namespace, kubeconfig string) *ApplyOptions {
	return &ApplyOptions{
		Dir:        dir,
		Suffix:     suffix,
		Namespace:  namespace,
		KubeConfig: kubeconfig,
	}
}

// Destroy executes terraform destroy.
func (c *Client) Destroy(o *ApplyOptions) error {
	tf, err := tfexec.NewTerraform(o.Dir, c.Binary)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	if err := tf.Init(context.TODO(),
		tfexec.Upgrade(true),
		tfexec.BackendConfig(fmt.Sprintf("secret_suffix=%s", o.Suffix)),
		tfexec.BackendConfig(fmt.Sprintf("namespace=%s", o.Namespace)),
		tfexec.BackendConfig(fmt.Sprintf("config_path=%s", o.KubeConfig))); err != nil {
		return err
	}

	if err := tf.Destroy(context.TODO()); err != nil {
		return err
	}
	return nil
}
