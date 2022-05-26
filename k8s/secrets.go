package k8s

import (
	"context"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (saCtx *CfgContext) ListSecrets() (*corev1.SecretList, error) {
	var ns string = ""

	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	return clientset.CoreV1().Secrets(ns).List(context.TODO(), metav1.ListOptions{})
}

func (saCtx *CfgContext) ListSecretNames() ([]string, error) {
	var ns string = ""
	var secretList = &corev1.SecretList{}
	var list []string = make([]string, 0)

	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	secretList, err = clientset.CoreV1().Secrets(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range secretList.Items {
		list = append(list, v.Name)
	}

	return list, nil
}

func (saCtx *CfgContext) GetSecret() (*corev1.Secret, error) {
	var ns string = ""
	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	return clientset.CoreV1().Secrets(ns).Get(context.TODO(), *saCtx.name, metav1.GetOptions{})
}

func (saCtx *CfgContext) CreateSecretFromObj(sa *corev1.Secret) (*corev1.Secret, error) {
	var ns string = ""
	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	return clientset.CoreV1().Secrets(ns).Create(context.TODO(), sa, metav1.CreateOptions{})
}

func (saCtx *CfgContext) CreateSecretFromYaml(yamlPath string) (*corev1.Secret, error) {
	var secret *corev1.Secret = &corev1.Secret{}
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlBytes, secret)
	if err != nil {
		return nil, err
	}

	return saCtx.CreateSecretFromObj(secret)
}
