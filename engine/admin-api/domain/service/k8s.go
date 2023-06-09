package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type K8sService interface {
	Start(ctx context.Context, productID string, version *entity.Version, versionConfig *entity.VersionConfig) error
	Stop(ctx context.Context, productID string, version *entity.Version) error
	Publish(ctx context.Context, productID string, version *entity.Version) error
	Unpublish(ctx context.Context, productID string, version *entity.Version) error
	WatchProcessStatus(ctx context.Context, productID, versionTag string) (<-chan *entity.Process, error)
}
