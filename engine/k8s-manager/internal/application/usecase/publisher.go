package usecase

import (
	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"golang.org/x/net/context"
)

type VersionPublisher struct {
	logger           logr.Logger
	networkPublisher service.NetworkPublisher
}

func NewVersionPublisher(logger logr.Logger, networkPublisher service.NetworkPublisher) VersionPublisherService {
	return &VersionPublisher{
		logger,
		networkPublisher,
	}
}

func (vp *VersionPublisher) PublishVersion(ctx context.Context, version domain.Version) error {
	return vp.networkPublisher.PublishNetwork(ctx, service.PublishNetworkParams{
		Product: version.Product,
		Version: version.Tag,
	})
}
