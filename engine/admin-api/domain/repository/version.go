package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type VersionRepo interface {
	Create(userEmail, productID string, version *entity.Version) (*entity.Version, error)
	CreateIndexes(ctx context.Context, productID string) error
	GetByTag(ctx context.Context, productID, tag string) (*entity.Version, error)
	GetLatest(ctx context.Context, productID string) (*entity.Version, error)
	ListVersionsByProduct(ctx context.Context, productID string, filter *ListVersionsFilter) ([]*entity.Version, error)
	Update(productID string, version *entity.Version) error
	// SetStatus updates the status and deletes the error message of the version.
	SetStatus(ctx context.Context, productID, versionTag string, status entity.VersionStatus) error
	SetErrorStatusWithError(ctx context.Context, productID, version, errorMessage string) error
	SetCriticalStatusWithError(ctx context.Context, productID, version, errorMessage string) error
}

type ListVersionsFilter struct {
	Status entity.VersionStatus
}
