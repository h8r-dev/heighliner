package state

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"github.com/h8r-dev/heighliner/pkg/state/app"
	"github.com/h8r-dev/heighliner/pkg/state/infra"
)

// ConfigMapState state using k8s configmap as backend
type ConfigMapState struct {
	ClientSet *kubernetes.Clientset
}

func (c *ConfigMapState) LoadInfra() (*infra.Output, error) {
	cm, err := c.ClientSet.CoreV1().ConfigMaps(InfraNs).Get(context.TODO(), InfraConfigMap, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	rawInfra := cm.Data[InfraEntry]
	infra := &infra.Output{}
	if err := yaml.Unmarshal([]byte(rawInfra), infra); err != nil {
		return nil, err
	}
	return infra, nil
}

// ListApps list all heighliner applications
func (c *ConfigMapState) ListApps() ([]string, error) {

	cms, err := c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.Set(map[string]string{configTypeKey: "heighliner"}).AsSelector().String(),
	})
	if err != nil {
		return nil, err
	}

	apps := make([]string, 0)
	for _, item := range cms.Items {
		apps = append(apps, item.Name)
	}
	return apps, nil
}

// LoadOutput load output from configmap
func (c *ConfigMapState) LoadOutput(appName string) (*app.Output, error) {

	cm, err := c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), appName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if cm.Data[stackOutput] == "" {
		return nil, fmt.Errorf("no data in configmap %s", appName)
	}

	ao := app.Output{}
	err = yaml.Unmarshal([]byte(cm.Data[stackOutput]), &ao)
	if err != nil {
		return nil, err
	}

	ao.ApplicationRef.Name = appName

	return &ao, nil
}

// LoadTFProvider Load tf provider from configmap
func (c *ConfigMapState) LoadTFProvider(appName string) (string, error) {

	cm, err := c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), appName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	tfConfigMapName := cm.Data[tfProviderConfigMapKey]
	if tfConfigMapName == "" {
		return "", fmt.Errorf("no tf provider config map? ")
	}

	cm, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), tfConfigMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if len(cm.Data) == 0 || cm.Data[tfProviderConfigMapKey] == "" {
		return "", fmt.Errorf("no data found in tf provider configmap")
	}
	return cm.Data[tfProviderConfigMapKey], nil
}

// SaveOutputAndTFProvider Save output and tf provider to configmap
func (c *ConfigMapState) SaveOutputAndTFProvider(appName string) error {
	ao, err := app.Load(stackOutput)
	if err != nil {
		return err
	}
	outputBys, err := os.ReadFile(stackOutput)
	if err != nil {
		return err
	}

	tfConfigName := "tf-" + appName
	configMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: appName, Labels: map[string]string{configTypeKey: "heighliner",
			"heighliner.dev/app-name": appName}},
		Data: map[string]string{stackOutput: string(outputBys), "tf-provider": tfConfigName},
	}

	// delete it if already exist
	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), appName, metav1.GetOptions{})
	if err == nil {
		_ = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Delete(context.TODO(), appName, metav1.DeleteOptions{})
	}

	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	tfBys, err := ioutil.ReadFile(ao.SCM.TfProvider)
	if err != nil {
		return fmt.Errorf("fail to read file from %s, err: %w", ao.SCM.TfProvider, err)
	}

	tfConfigMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: tfConfigName, Labels: map[string]string{configTypeKey: "tf-provider",
			"heighliner.dev/app-name": appName}},
		Data: map[string]string{tfProviderConfigMapKey: string(tfBys)},
	}

	// delete it if already exist
	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), tfConfigName, metav1.GetOptions{})
	if err == nil {
		_ = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Delete(context.TODO(), tfConfigName, metav1.DeleteOptions{})
	}

	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Create(context.TODO(), &tfConfigMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return os.Remove(stackOutput)
}

// DeleteOutputAndTFProvider delete output and tf provider configMap
func (c *ConfigMapState) DeleteOutputAndTFProvider(appName string) error {
	ctx := context.TODO()
	tfConfigName := "tf-" + appName
	if err := c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Delete(ctx, appName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Delete(ctx, tfConfigName, metav1.DeleteOptions{})
}
