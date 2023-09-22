package service

import (
	"context"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
)

type CreateProcessParams struct {
	ConfigName string
	Product    string
	Version    string
	Workflow   string
	Process    *domain.Process
}

type CreateNetworkParams struct {
	Product  string
	Version  string
	Workflow string
	Process  *domain.Process
}

type PublishNetworkParams struct {
	Product string
	Version string
}

type ContainerStarter interface {
	CreateProcess(ctx context.Context, params CreateProcessParams) error
	CreateNetwork(ctx context.Context, params CreateNetworkParams) error
	CreateVersionConfiguration(ctx context.Context, version domain.Version) (string, error)
}

type ContainerStopper interface {
	DeleteProcesses(ctx context.Context, product, version string) error
	DeleteConfiguration(ctx context.Context, product, version string) error
	DeleteNetwork(ctx context.Context, product, version string) error
}

type ContainerPublisher interface {
	PublishNetwork(ctx context.Context, params PublishNetworkParams) (map[string]string, error)
}

//go:generate mockery --name ImageBuilder --output ../../../mocks --filename image_builder_mock.go --structname ImageBuilderMock
type ImageBuilder interface {
	BuildImage(ctx context.Context, processID, processImage string, sources []byte) (string, error)
}

//go:generate mockery --name ContainerService --output ../../../mocks --filename container_service_mock.go --structname ContainerServiceMock
type ContainerService interface {
	ContainerStarter
	ContainerStopper
	ContainerPublisher
}
