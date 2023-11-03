package usecase

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
)

//go:generate mockery --name ProcessService --output ../../../mocks --filename process_service_mock.go --structname ProcessServiceMock

type ProcessService interface {
	RegisterProcess(ctx context.Context, params RegisterProcessParams) (string, error)
}

type ProcessRegister struct {
	logger       logr.Logger
	imageBuilder service.ImageBuilder
}

func NewProcessRegister(logger logr.Logger, imageRegistry service.ImageBuilder) *ProcessRegister {
	return &ProcessRegister{
		logger:       logger,
		imageBuilder: imageRegistry,
	}
}

type RegisterProcessParams struct {
	ProductID    string
	ProcessID    string
	ProcessImage string
}

func (pr *ProcessRegister) RegisterProcess(ctx context.Context, params RegisterProcessParams) (string, error) {
	imageID, err := pr.imageBuilder.BuildImage(ctx, params.ProductID, params.ProcessID, params.ProcessImage)
	if err != nil {
		return "", fmt.Errorf("building image: %w", err)
	}

	return imageID, nil
}
