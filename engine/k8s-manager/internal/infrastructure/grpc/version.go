package grpc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/infrastructure/grpc/proto/versionpb"
)

type VersionService struct {
	versionpb.UnimplementedVersionServiceServer
	logger          logr.Logger
	starter         usecase.VersionStarterService
	stopper         usecase.VersionStopperService
	publisher       usecase.VersionPublisherService
	unpublisher     usecase.VersionUnpublisherService
	processRegister usecase.ProcessService
}

func NewVersionService(
	logger logr.Logger,
	starter usecase.VersionStarterService,
	stopper usecase.VersionStopperService,
	publisher usecase.VersionPublisherService,
	unpublisher usecase.VersionUnpublisherService,
	processRegister usecase.ProcessService,
) *VersionService {
	return &VersionService{
		versionpb.UnimplementedVersionServiceServer{},
		logger,
		starter,
		stopper,
		publisher,
		unpublisher,
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
) (*versionpb.RegisterProcessResponse, error) {
	v.logger.Info("Register process request received")

	imageID, err := v.processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		ProcessID:    req.ProcessId,
		ProcessImage: req.ProcessImage,
		Sources:      req.File,
	})
	if err != nil {
		return nil, fmt.Errorf("registering process: %w", err)
	}

	return &versionpb.RegisterProcessResponse{
		ImageId: imageID,
	}, nil
}

func (v *VersionService) Publish(
	ctx context.Context,
	req *versionpb.PublishRequest,
) (*versionpb.PublishResponse, error) {
	v.logger.Info("Register process request received")

	networkURLs, err := v.publisher.PublishVersion(ctx, req.Product, req.VersionTag)
	if err != nil {
		return nil, fmt.Errorf("registering process: %w", err)
	}

	fmt.Println(networkURLs)

	return &versionpb.PublishResponse{
		NetworkUrls: networkURLs,
	}, nil
}

func (v *VersionService) Unpublish(
	ctx context.Context,
	req *versionpb.UnpublishRequest,
) (*versionpb.Response, error) {
	v.logger.Info("Unpublish request received")

	err := v.unpublisher.UnpublishVersion(ctx, req.Product, req.VersionTag)
	if err != nil {
		return nil, fmt.Errorf("unpublishing version: %w", err)
	}

	return &versionpb.Response{
		Message: fmt.Sprintf("Version %q on product %q unpublished", req.VersionTag, req.Product),
	}, nil
}
