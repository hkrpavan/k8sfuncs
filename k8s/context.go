package k8s

type CfgContext struct {
	clusterName, cfgPath, namespace, name, user *string
}

func CreateK8SContext(clusterName, cfgPath, namespace, name, user string) CfgContext {
	return CfgContext{
		clusterName: &clusterName,
		cfgPath:     &cfgPath,
		namespace:   &namespace,
		name:        &name,
		user:        &user,
	}
}
