package grpc

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"

	"github.com/go-logr/logr"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
)

type VersionService struct {
	versionpb.UnimplementedVersionServiceServer
	logger  logr.Logger
	starter usecase.VersionStarterService
	stopper usecase.VersionStopperService
}

func NewVersionService(
	logger logr.Logger,
	starter usecase.VersionStarterService,
	stopper usecase.VersionStopperService,
) *VersionService {
	return &VersionService{
		versionpb.UnimplementedVersionServiceServer{},
		logger,
		starter,
		stopper,
	}
}

func (v *VersionService) Start(
	ctx context.Context, req *versionpb.StartRequest,
) (*versionpb.Response, error) {
	v.logger.Info("Start request received")

	err := v.starter.StartVersion(ctx, mapRequestToVersion(req))
	if err != nil {
		return nil, fmt.Errorf("start version %q in product %q: %w", req.VersionTag, req.ProductId, err)
	}

	return &versionpb.Response{
		Message: fmt.Sprintf("Version %q in product %q started", req.VersionTag, req.ProductId),
	}, nil
}

func (v *VersionService) Stop(
	ctx context.Context,
	req *versionpb.StopRequest,
) (*versionpb.Response, error) {
	v.logger.Info("Stop request received")

	err := v.stopper.StopVersion(ctx, usecase.StopParams{
		Product: req.Product,
		Version: req.VersionTag,
	})
	if err != nil {
		return nil, fmt.Errorf("stop version %q in product %q: %w", req.VersionTag, req.Product, err)
	}

	return &versionpb.Response{
		Message: fmt.Sprintf("Version %q on product %q stopped", req.VersionTag, req.Product)}, nil
}
