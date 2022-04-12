package k8s

import (
	"context"
	"time"

	"golang.cisco.com/accordion/a7nlibs/logging"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type CfgContext struct {
	masterUrl, cfgPath, namespace, name *string
}

func CreatePodContext(masterUrl, cfgPath, namespace, name string) CfgContext {
	return CfgContext{
		masterUrl: &masterUrl,
		cfgPath:   &cfgPath,
		namespace: &namespace,
		name:      &name,
	}
}

func (podCtx *CfgContext) ListPods() []corev1.Pod {
	var cfgPath string
	var podLIst *corev1.PodList
	var cfg *rest.Config = GetUserK8SConfig(podCtx)
	var clientset, err = kubernetes.NewForConfig(cfg)
	var ns string = ""

	if err != nil {
		logging.ErrorLogger.Printf("Error !! Creating clientset for the config, %s.\n", cfgPath)
		return nil
	}

	if *podCtx.namespace != "" {
		ns = *podCtx.namespace
	}

	podLIst, err = clientset.CoreV1().Pods(ns).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		logging.ErrorLogger.Println("Error !! Unable to retieve pod list for given context.")
		return nil
	}

	return podLIst.Items

}

func ListRunningPodNames(list []corev1.Pod) *[]string {
	pods := make([]string, 0)
	for _, v := range list {
		if v.Status.Phase == "Running" {
			pods = append(pods, v.Name)
		}
	}
	return &pods
}

func ListFailedPodNames(list []corev1.Pod) *[]string {
	pods := make([]string, 0)
	for _, v := range list {
		if v.Status.Phase == "Failed" {
			pods = append(pods, v.Name)
		}
	}
	return &pods
}

func ListEvictedPodNames(list []corev1.Pod) *[]string {
	pods := make([]string, 0)
	for _, v := range list {
		if v.Status.Phase == "Evicted" {
			pods = append(pods, v.Name)
		}
	}
	return &pods
}

func ListPendingPodNames(list []corev1.Pod) *[]string {
	pods := make([]string, 0)
	for _, v := range list {
		if v.Status.Phase == "Pending" {
			pods = append(pods, v.Name)
		}
	}
	return &pods
}

func ListPodNames(list []corev1.Pod) *[]string {
	pods := make([]string, 0)
	for _, v := range list {
		pods = append(pods, v.Name)
	}
	return &pods
}

func (podCtx *CfgContext) CreatePodFromSpecObj(podSpeclObj *corev1.Pod) (*corev1.Pod, error) {
	var cfg *rest.Config = GetUserK8SConfig(podCtx)
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logging.ErrorLogger.Printf("Error !! Unable to create clientset for given kube config %s.\n", *podCtx.cfgPath)
		panic(err.Error())
	}

	return clientset.CoreV1().Pods(podSpeclObj.Namespace).Create(context.TODO(), podSpeclObj, v1.CreateOptions{})

}

func (podCtx *CfgContext) GetPodStatus() (*corev1.Pod, error) {
	var cfg *rest.Config = GetUserK8SConfig(podCtx)
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		logging.ErrorLogger.Printf("Error !! Unable to create clientset for given kube config %s.\n", *podCtx.cfgPath)
		panic(err.Error())
	}

	return clientset.CoreV1().Pods(*podCtx.namespace).Get(context.TODO(), *podCtx.name, v1.GetOptions{})
}

func (podCtx *CfgContext) GetDeployedPodStatusRetry(count int8) *corev1.Pod {
	var i int8 = 1
	var status *corev1.Pod
	var seconds int8 = 2
	for i <= count {
		status, _ = podCtx.GetPodStatus()
		if i > 1 {
			logging.InfoLogger.Printf("Status check for Pod \"%s\" in namespace \"%s\". Retry %d\n", status.Name, status.Namespace, i)
		}
		if status.Status.Phase == "Running" {
			logging.InfoLogger.Printf("Pod \"%s\" is in Running status in namespace: \"%s\"\n", status.Name, status.Namespace)
			break
		} else if status.Status.Phase == "Pending" {
			logging.InfoLogger.Printf("Pod \"%s\" is still in Pending status in namespace: \"%s\"\n", status.Name, status.Namespace)
			logging.InfoLogger.Printf("Sleeping for %d seconds..", seconds)
			time.Sleep(time.Duration(seconds) * time.Second)
			i++
			continue
		} else {
			logging.ErrorLogger.Printf("Pod creation failed for %s in namespace %s", status.Name, status.Namespace)
			break
		}
	}
	if status.Status.Phase != "Running" {
		logging.ErrorLogger.Printf("Pod not in Running state  with in %d seconds. Retried for %d times, aborting.", seconds*count, count)
	}
	return status
}
