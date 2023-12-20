package configuration

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

type KubeConfiguration struct {
	logger    logr.Logger
	client    kubernetes.Interface
	namespace string
}

func NewKubeConfiguration(logger logr.Logger, client kubernetes.Interface, namespace string) KubeConfiguration {
	return KubeConfiguration{
		logger:    logger,
		client:    client,
		namespace: namespace,
	}
}
