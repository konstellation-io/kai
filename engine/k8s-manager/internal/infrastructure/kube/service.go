package kube

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/configuration"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/network"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/kube/process"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
)

type K8sContainerService struct {
	logger    logr.Logger
	namespace string
	client    kubernetes.Interface

	processService       process.KubeProcess
	configurationService configuration.KubeConfiguration
	networkService       network.KubeNetwork
}

var _ service.ContainerService = (*K8sContainerService)(nil)

func NewK8sContainerService(logger logr.Logger, client kubernetes.Interface) *K8sContainerService {
	namespace := viper.GetString("kubernetes.namespace")

	return &K8sContainerService{
		logger:    logger,
		namespace: namespace,
		client:    client,

		processService:       process.NewKubeProcess(logger, client, namespace),
		configurationService: configuration.NewKubeConfiguration(logger, client, namespace),
		networkService:       network.NewKubeNetwork(logger, client, namespace),
	}
}

func (k *K8sContainerService) CreateProcess(ctx context.Context, params service.CreateProcessParams) error {
	return k.processService.Create(ctx, params)
}

func (k *K8sContainerService) DeleteProcesses(ctx context.Context, product, version string) error {
	return k.processService.DeleteProcesses(ctx, product, version)
}

func (k *K8sContainerService) CreateVersionConfiguration(ctx context.Context, version domain.Version) (string, error) {
	return k.configurationService.CreateVersionConfiguration(ctx, version)
}
func (k *K8sContainerService) DeleteConfiguration(ctx context.Context, product, version string) error {
	return k.configurationService.DeleteConfiguration(ctx, product, version)
}

func (k *K8sContainerService) CreateNetwork(ctx context.Context, params service.CreateNetworkParams) error {
	return k.networkService.CreateNetwork(ctx, params)
}

func (k *K8sContainerService) DeleteNetwork(ctx context.Context, product, version string) error {
	return k.networkService.DeleteNetwork(ctx, product, version)
}
