package usecase

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"golang.org/x/net/context"
)

type VersionPublisher struct {
	logger           logr.Logger
	networkPublisher service.ContainerPublisher
}

func NewVersionPublisher(logger logr.Logger, networkPublisher service.ContainerPublisher) VersionPublisherService {
	return &VersionPublisher{
		logger,
		networkPublisher,
	}
}

func (vp *VersionPublisher) PublishVersion(ctx context.Context, product, version string) (map[string]string, error) {
	vp.logger.Info("Publishing version", "product", product, "version", version)

	return vp.networkPublisher.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	})
}

func (vp *VersionPublisher) GetPublishedTriggers(ctx context.Context, product string) (map[string]string, error) {
	vp.logger.Info("Getting published triggers", "product", product)

	return vp.networkPublisher.GetPublishedTriggers(ctx, product)
}
