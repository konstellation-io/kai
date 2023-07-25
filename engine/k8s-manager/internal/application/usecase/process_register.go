package usecase

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
)

type ProcessRegister struct {
	logger        logr.Logger
	imageRegistry service.ImageRegistry
}

func NewProcessRegister(logger logr.Logger, imageRegistry service.ImageRegistry) *ProcessRegister {
	return &ProcessRegister{
		logger:        logger,
		imageRegistry: imageRegistry,
	}
}

type RegisterProcessParams struct {
	Product string
	Version string
	Process string
	File    []byte
}

func (pr *ProcessRegister) RegisterProcess(ctx context.Context, params RegisterProcessParams) error {
	imageName := fmt.Sprintf("kai-local-registry:5000/%s-%s:%s", params.Product, params.Process, params.Version)

	return pr.imageRegistry.ProcessRegister(ctx, imageName, params.File)
}
