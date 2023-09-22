package usecase

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
)

type VersionUnpublisher struct {
	logger             logr.Logger
	networkUnpublisher service.ContainerUnpublisher
}

func NewVersionUnpublisher(logger logr.Logger, networkUnpublisher service.ContainerUnpublisher) VersionUnpublisherService {
	return &VersionUnpublisher{
		logger,
		networkUnpublisher,
	}
}

func (vp *VersionUnpublisher) UnpublishVersion(ctx context.Context, product, version string) error {
	return vp.networkUnpublisher.UnpublishNetwork(ctx, product, version)
}
