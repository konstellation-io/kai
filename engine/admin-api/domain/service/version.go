package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type VersionService interface {
	Start(ctx context.Context, product *entity.Product, version *entity.Version, versionConfig *entity.VersionStreamingResources) error
	Stop(ctx context.Context, productID string, version *entity.Version) error
	Publish(ctx context.Context, productID, versionTag string) (map[string]string, error)
	Unpublish(ctx context.Context, productID string, version *entity.Version) error
	WatchProcessStatus(ctx context.Context, productID, versionTag string) (<-chan *entity.Process, error)
	RegisterProcess(ctx context.Context, productID, processID, processImage string) (string, error)
}
