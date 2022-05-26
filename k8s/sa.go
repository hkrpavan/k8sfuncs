package k8s

import (
	"context"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (saCtx *CfgContext) ListSAs() ([]corev1.ServiceAccount, error) {
	var saLIst *corev1.ServiceAccountList
	var ns string = ""

	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	saLIst, err = clientset.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return saLIst.Items, nil
}

func (saCtx *CfgContext) ListSANames() ([]string, error) {
	var saLIst *corev1.ServiceAccountList
	var ns string = ""
	var saNames []string = make([]string, 0)

	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	saLIst, err = clientset.CoreV1().ServiceAccounts(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, v := range saLIst.Items {
		saNames = append(saNames, v.Name)
	}

	return saNames, nil
}

func (saCtx *CfgContext) GetServiceAccount() (*corev1.ServiceAccount, error) {
	var ns string = ""
	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	return clientset.CoreV1().ServiceAccounts(ns).Get(context.TODO(), *saCtx.name, metav1.GetOptions{})
}

func (saCtx *CfgContext) GetServiceAccountSecrets(name *string) ([]string, error) {
	var secrets []string = make([]string, 0)

	sa, err := saCtx.GetServiceAccount()
	if err != nil {
		return nil, err
	}

	for _, v := range sa.Secrets {
		secrets = append(secrets, v.Name)
	}
	return secrets, nil
}

func (saCtx *CfgContext) CreateServiceAccount(name *string) (*corev1.ServiceAccount, error) {
	var sa *corev1.ServiceAccount = &corev1.ServiceAccount{}
	var ns string
	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	sa.Name = *name

	return clientset.CoreV1().ServiceAccounts(ns).Create(context.TODO(), sa, metav1.CreateOptions{})
}

func (saCtx *CfgContext) CreateSAFromObj(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	var ns string
	clientset, err := saCtx.GetClientForConfig()
	if err != nil {
		return nil, err
	}

	if *saCtx.namespace != "" {
		ns = *saCtx.namespace
	}

	return clientset.CoreV1().ServiceAccounts(ns).Create(context.TODO(), sa, metav1.CreateOptions{})
}

func (saCtx *CfgContext) CreateSAFromYaml(yamlPath string) (*corev1.ServiceAccount, error) {
	var sa *corev1.ServiceAccount = &corev1.ServiceAccount{}
	yamlBytes, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlBytes, sa)
	if err != nil {
		return nil, err
	}

	return saCtx.CreateSAFromObj(sa)
}
