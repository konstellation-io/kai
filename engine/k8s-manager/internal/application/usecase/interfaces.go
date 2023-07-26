package usecase

import (
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	"golang.org/x/net/context"
)

type VersionStarterService interface {
	StartVersion(ctx context.Context, version domain.Version) error
}

type VersionStopperService interface {
	StopVersion(ctx context.Context, params StopParams) error
}

//go:generate mockgen -source=${GOFILE} -destination=../../../mocks/version_service_mock.go -package=mocks
type VersionService interface {
	VersionStarterService
	VersionStopperService
}
