package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProcessRepository interface {
	CreateIndexes(ctx context.Context, productID string) error
	Create(productID string, newRegisteredProcess *entity.RegisteredProcess) (*entity.RegisteredProcess, error)
	ListByProductAndType(ctx context.Context, productID, processType string) ([]*entity.RegisteredProcess, error)
}
