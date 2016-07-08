package kubeclient

import (
	"github.com/golang/glog"
        "k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"

        kube_restclient "k8s.io/kubernetes/pkg/client/restclient"
	kube_client "k8s.io/kubernetes/pkg/client/unversioned"
	kube_clientcmd "k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	kube_clientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
)

const (
	APIVersion = "v1"
)
/*
var (
	kubeClient *kube_client.Client
	KubeConfig = &Config{
		APIServerURL: "https://61.160.36.122",
		Username:     "test",
		Password:     "test123",
	}
)

type Config struct {
	APIServerURL string
	Username     string
	Password     string
}

func Init() {
	var err error
	kubeClient, err = getKubeClient()
	if err != nil {
		glog.Fatalf("Can not connect to kubernetes: %v", err)
	}
}

func Get() *kube_client.Client {
	if kubeClient == nil {
		glog.Fatalf("Forget to call kubeclient.Init()?")
	}
	return kubeClient
}

func getConfigOverrides() (*kube_clientcmd.ConfigOverrides, error) {
	kubeConfigOverride := kube_clientcmd.ConfigOverrides{
		ClusterInfo: kube_clientcmdapi.Cluster{
			APIVersion: APIVersion,
		},
	}

	kubeConfigOverride.ClusterInfo.Server = KubeConfig.APIServerURL
	kubeConfigOverride.ClusterInfo.InsecureSkipTLSVerify = true

	return &kubeConfigOverride, nil
}

func getKubeClient() (*kube_client.Client, error) {
	configOverrides, err := getConfigOverrides()
	if err != nil {
		return nil, err
	}

	kubeConfig := &kube_restclient.Config{
		Host:     configOverrides.ClusterInfo.Server,
		Insecure: configOverrides.ClusterInfo.InsecureSkipTLSVerify,
		Username: KubeConfig.Username,
		Password: KubeConfig.Password,
	}

        kubeConfig.GroupVersion = &unversioned.GroupVersion{Version: configOverrides.ClusterInfo.APIVersion}


	kubeClient := kube_client.NewOrDie(kubeConfig)
	return kubeClient, nil
}

func GetAllPods() ([]*api.Pod, error) {

        opt := api.ListOptions {
	LabelSelector: labels.Everything(),
	FieldSelector: fields.Everything(),
        }

	podList, err := kubeClient.Pods("").List(opt)
	if err != nil {
		return nil, err
	}
	var result []*api.Pod
	for j := range podList.Items {
		pod := &podList.Items[j]
		result = append(result, pod)
	}
	return result, nil
}

*/
