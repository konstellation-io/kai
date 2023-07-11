package network

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

type KubeNetwork struct {
	logger    logr.Logger
	client    kubernetes.Interface
	namespace string
}

func NewKubeNetwork(logger logr.Logger, client kubernetes.Interface, namespace string) KubeNetwork {
	return KubeNetwork{
		logger:    logger,
		client:    client,
		namespace: namespace,
	}
}
