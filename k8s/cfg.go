package k8s

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"golang.cisco.com/accordion/a7nlibs/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func validateKubeCfgFile(cfgPath *string) error {
	cfg := clientcmd.GetConfigFromFileOrDie(*cfgPath)
	if (len(reflect.ValueOf(cfg.Clusters).MapKeys())) > 1 {
		logging.InfoLogger.Printf("Found more than one cluster in \"%s\"", *cfgPath)
		logging.InfoLogger.Println("Checking for current context in kubeconfig file..")
		logging.InfoLogger.Printf("Current context is set to %s.\n", cfg.CurrentContext)
		if cfg.CurrentContext == "" {
			msg := "missing current context in kubeconfig file"
			return errors.New(msg)
		}
	}
	return nil
}

func (ctxCfg *CfgContext) GetUserK8SConfig() (*rest.Config, error) {
	var err error
	var cfg *rest.Config
	var usercfgPath string

	if *ctxCfg.cfgPath == "" {
		userHome := homedir.HomeDir()
		usercfgPath = fmt.Sprintf("%s/.kube/config", userHome)
		_, e := os.Stat(usercfgPath)
		if os.IsNotExist(e) {
			logging.ErrorLogger.Printf("Error !! Kuber config file not found in default location %s.\n", usercfgPath)
			return nil, e
		} else {
			*ctxCfg.cfgPath = usercfgPath
		}
	}

	if err = validateKubeCfgFile(ctxCfg.cfgPath); err != nil {
		return nil, err
	}
	logging.InfoLogger.Printf("Building kubeconfig for file %s\n", *ctxCfg.cfgPath)
	cfg, err = clientcmd.BuildConfigFromFlags(*ctxCfg.clusterName, *ctxCfg.cfgPath)
	if err != nil {
		logging.ErrorLogger.Printf("Unable to fetch kubeconfig for %s from %s.\n", *ctxCfg.clusterName, *ctxCfg.cfgPath)
		return nil, err
	}

	return cfg, nil
}

func (ctxCfg *CfgContext) GetClientForConfig() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	config, err := ctxCfg.GetUserK8SConfig()
	if err != nil {
		return nil, err
	}
	clientset, err = kubernetes.NewForConfig(config)
	return clientset, err
}

func (ctxCfg *CfgContext) NewKubeConfig() (*clientcmdapi.Config, error) {
	//var cfg *rest.Config
	var err error
	var kubectlcfg *clientcmdapi.Config = &clientcmdapi.Config{
		APIVersion: "v1",
		Kind:       "Config",
	}
	//cfg, err = ctxCfg.GetUserK8SConfig()
	if err != nil {
		return nil, err
	}
	return kubectlcfg, nil
}
