package service

import (
	"context"

	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
)

//go:generate mockery --name ContainerService

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

type ImageRegistry interface {
	ProcessRegister(ctx context.Context, name string, file []byte) error
}

type ContainerService interface {
	ContainerStarter
	ContainerStopper
	ImageRegistry
}
