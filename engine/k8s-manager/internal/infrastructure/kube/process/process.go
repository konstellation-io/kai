package process

import (
	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
)

type KubeProcess struct {
	logger    logr.Logger
	client    kubernetes.Interface
	namespace string
}

func NewKubeProcess(logger logr.Logger, client kubernetes.Interface, namespace string) KubeProcess {
	return KubeProcess{
		logger:    logger,
		client:    client,
		namespace: namespace,
	}
}
