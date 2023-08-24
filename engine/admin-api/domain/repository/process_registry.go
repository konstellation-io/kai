package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProcessRegistryRepo interface {
	CreateIndexes(ctx context.Context, productID string) error
	Create(productID string, newProcessRegistry *entity.ProcessRegistry) (*entity.ProcessRegistry, error)
	ListByProductWithTypeFilter(ctx context.Context, productID, processType string) ([]*entity.ProcessRegistry, error)
}
