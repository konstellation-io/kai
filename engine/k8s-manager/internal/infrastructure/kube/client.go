package kube

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/config"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClientset() (kubernetes.Interface, error) {
	kubeConfig, err := newKubernetesConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func newKubernetesConfig() (*rest.Config, error) {
	if viper.GetBool("kubernetes.isInsideCluster") {
		kubeConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		return kubeConfig, nil
	}

	// use the current context in kubeConfigPath
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", viper.GetString(config.KubeConfigPathKey))
	if err != nil {
		return nil, err
	}

	return kubeConfig, nil
}
