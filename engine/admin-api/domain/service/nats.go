package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type NatsManagerService interface {
	CreateStreams(ctx context.Context, product string, version *entity.Version) (*entity.VersionStreamsConfig, error)
	CreateObjectStores(ctx context.Context, product string, version *entity.Version) (*entity.VersionObjectStoresConfig, error)
	DeleteStreams(ctx context.Context, product string, versionTag string) error
	DeleteObjectStores(ctx context.Context, product, versionTag string) error
	CreateVersionKeyValueStores(ctx context.Context, product string, version *entity.Version) (*entity.KeyValueStores, error)
	CreateGlobalKeyValueStore(ctx context.Context, product string) (string, error)
	UpdateKeyValueConfiguration(ctx context.Context, configurations []entity.KeyValueConfiguration) error
}
