package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type VersionService interface {
	Start(ctx context.Context, runtimeID string, version *entity.Version, versionConfig *entity.VersionConfig) error
	Stop(ctx context.Context, runtimeID string, version *entity.Version) error
	Publish(runtimeID string, version *entity.Version) error
	Unpublish(runtimeID string, version *entity.Version) error
	UpdateConfig(runtimeID string, version *entity.Version) error
	WatchProcessStatus(ctx context.Context, runtimeID, versionName string) (<-chan *krt.Process, error)
}
