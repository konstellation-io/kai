package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProcessRepository interface {
	CreateIndexes(ctx context.Context, registry string) error
	Create(ctx context.Context, registry string, newRegisteredProcess *entity.RegisteredProcess) error
	SearchByProduct(ctx context.Context, product string, filter *entity.SearchFilter) ([]*entity.RegisteredProcess, error)
	GlobalSearch(ctx context.Context, filter *entity.SearchFilter) ([]*entity.RegisteredProcess, error)
	Update(ctx context.Context, registry string, newRegisteredProcess *entity.RegisteredProcess) error
	GetByID(ctx context.Context, registry string, imageID string) (*entity.RegisteredProcess, error)
	Delete(ctx context.Context, registry string, imageID string) error
}
