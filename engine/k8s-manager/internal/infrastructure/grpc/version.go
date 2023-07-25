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
	logger          logr.Logger
	starter         *usecase.VersionStarter
	stopper         *usecase.VersionStopper
	processRegister *usecase.ProcessRegister
}

func NewVersionService(
	logger logr.Logger,
	starter *usecase.VersionStarter,
	stopper *usecase.VersionStopper,
	processRegister *usecase.ProcessRegister,
) *VersionService {
	return &VersionService{
		versionpb.UnimplementedVersionServiceServer{},
		logger,
		starter,
		stopper,
		processRegister,
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

func (v *VersionService) RegisterProcess(
	ctx context.Context,
	req *versionpb.RegisterProcessRequest,
) (*versionpb.Response, error) {
	v.logger.Info("Register process request received")

	err := v.processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		Product: req.Product,
		Version: req.Version,
		Process: req.Process,
		File:    req.File,
	})
	if err != nil {
		return nil, fmt.Errorf("registering process: %w", err)
	}

	return &versionpb.Response{
		Message: fmt.Sprintf("Registered process")}, nil
}
