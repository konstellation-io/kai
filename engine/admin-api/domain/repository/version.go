package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type VersionRepo interface {
	Create(userEmail, productID string, version *entity.Version) (*entity.Version, error)
	CreateIndexes(ctx context.Context, productID string) error
	GetByID(productID, versionID string) (*entity.Version, error)
	GetByTag(ctx context.Context, productID, tag string) (*entity.Version, error)
	ListVersionsByProduct(ctx context.Context, productID string) ([]*entity.Version, error)
	Update(productID string, version *entity.Version) error
	ClearPublishedVersion(ctx context.Context, productID string) (*entity.Version, error)

	// SetStatus updates the status and deletes the error message of the version.
	SetStatus(ctx context.Context, productID, versionTag string, status entity.VersionStatus) error
	// SetError sets the error message of the version and updates the status to Error.
	SetError(ctx context.Context, productID string, version *entity.Version, errorMessage string) (*entity.Version, error)
}
