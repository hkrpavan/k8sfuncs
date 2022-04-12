package k8s

import (
	"fmt"
	"os"

	"golang.cisco.com/accordion/a7nlibs/logging"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetUserK8SConfig(podCtx *CfgContext) *rest.Config {
	var err error
	var cfg *rest.Config

	if *podCtx.cfgPath == "" {
		userHome := homedir.HomeDir()
		usercfgPath := fmt.Sprintf("%s/.kube/config", userHome)
		_, e := os.Stat(usercfgPath)
		if os.IsNotExist(e) {
			logging.ErrorLogger.Printf("Error !! Kuber config file not found in default location %s.\n", usercfgPath)
			panic(err.Error())
		} else {
			*podCtx.cfgPath = usercfgPath
		}
	}

	cfg, err = clientcmd.BuildConfigFromFlags(*podCtx.masterUrl, *podCtx.cfgPath)
	if err != nil {
		logging.ErrorLogger.Printf("Error !! retrieving kubernetes config from file %s.\n", *podCtx.cfgPath)
		panic(err.Error())
	}

	return cfg
}
