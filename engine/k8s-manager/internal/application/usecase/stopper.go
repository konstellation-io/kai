package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
)

type StopParams struct {
	Product string
	Version string
}

type VersionStopper struct {
	logger           logr.Logger
	containerService service.ContainerStopper
}

func NewVersionStopper(logger logr.Logger, containerService service.ContainerStopper) *VersionStopper {
	return &VersionStopper{
		logger:           logger,
		containerService: containerService,
	}
}

func (s *VersionStopper) StopVersion(ctx context.Context, params StopParams) error {
	product := params.Product
	version := params.Version

	s.logger.Info("Stopping version", "product", product, "version", version)

	var errs error
	if err := s.containerService.DeleteConfiguration(ctx, product, version); err != nil {
		errs = errors.Join(errs, fmt.Errorf("delete configuration: %w", err))
	}

	if err := s.containerService.DeleteNetwork(ctx, product, version); err != nil {
		errs = errors.Join(errs, fmt.Errorf("delete network: %w", err))
	}

	if err := s.containerService.DeleteProcesses(ctx, product, version); err != nil {
		errs = errors.Join(errs, fmt.Errorf("delete processes: %w", err))
	}

	return errs
}
