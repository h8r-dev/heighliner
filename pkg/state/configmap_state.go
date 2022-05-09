package state

import (
	"context"
	"fmt"
	"github.com/h8r-dev/heighliner/pkg/state/app"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

type ConfigMapState struct {
	ClientSet *kubernetes.Clientset
}

func (c *ConfigMapState) LoadOutput(appName string) (*app.Output, error) {

	cm, err := c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Get(context.TODO(), appName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	//	Get(context.TODO(), metav1.ListOptions{
	//	LabelSelector: labels.Set(map[string]string{ConfigTypeKey: "heighliner"}).AsSelector().String(),
	//})

	if cm.Data["output.yaml"] == "" {
		return nil, fmt.Errorf("no data in configmap %s", appName)
	}

	ao := app.Output{}
	err = yaml.Unmarshal([]byte(cm.Data["output.yaml"]), &ao)
	if err != nil {
		return nil, err
	}

	return &ao, nil
}

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

func (c *ConfigMapState) SaveOutputAndTFProvider(appName string) error {
	outputBys, err := ioutil.ReadFile(stackOutput)
	if err != nil {
		return err
	}

	tfConfigName := "tf-" + appName
	configMap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: appName, Labels: map[string]string{configTypeKey: "heighliner",
			"heighliner.dev/app-name": appName}},
		Data: map[string]string{"output.yaml": string(outputBys), "tf-provider": tfConfigName},
	}

	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Create(context.TODO(), &configMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	ao := app.Output{}
	if err = yaml.Unmarshal(outputBys, &ao); err != nil {
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
	_, err = c.ClientSet.CoreV1().ConfigMaps(HeighlinerNs).Create(context.TODO(), &tfConfigMap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return os.Remove(stackOutput)
}
