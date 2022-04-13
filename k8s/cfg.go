package k8s

import (
	"fmt"
	"os"

	"golang.cisco.com/accordion/a7nlibs/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func (ctxCfg *CfgContext) GetUserK8SConfig() (*rest.Config, error) {
	var err error
	var cfg *rest.Config

	if *ctxCfg.cfgPath == "" {
		userHome := homedir.HomeDir()
		usercfgPath := fmt.Sprintf("%s/.kube/config", userHome)
		_, e := os.Stat(usercfgPath)
		if os.IsNotExist(e) {
			logging.ErrorLogger.Printf("Error !! Kuber config file not found in default location %s.\n", usercfgPath)
			return nil, e
		} else {
			*ctxCfg.cfgPath = usercfgPath
		}
	}

	cfg, err = clientcmd.BuildConfigFromFlags(*ctxCfg.masterUrl, *ctxCfg.cfgPath)
	if err != nil {
		logging.ErrorLogger.Printf("Error !! retrieving kubernetes config from file %s.\n", *ctxCfg.cfgPath)
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
