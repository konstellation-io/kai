package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type VersionRepo interface {
	Create(userID, productID string, version *entity.Version) (*entity.Version, error)
	CreateIndexes(ctx context.Context, productID string) error
	GetByID(productID, versionID string) (*entity.Version, error)
	GetByTag(ctx context.Context, productID, tag string) (*entity.Version, error)
	ListVersionsByProduct(ctx context.Context, productID string) ([]*entity.Version, error)
	Update(productID string, version *entity.Version) error
	SetStatus(ctx context.Context, productID, versionID string, status entity.VersionStatus) error
	SetError(ctx context.Context, productID string, version *entity.Version, errorMessage string) (*entity.Version, error)
	UploadKRTYamlFile(productID string, version *entity.Version, file string) error
	ClearPublishedVersion(ctx context.Context, productID string) (*entity.Version, error)
}
