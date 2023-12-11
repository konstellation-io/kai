package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProcessRepository interface {
	CreateIndexes(ctx context.Context, productID string) error
	Create(productID string, newRegisteredProcess *entity.RegisteredProcess) (*entity.RegisteredProcess, error)
	SearchByProduct(ctx context.Context, product string, filter SearchFilter) ([]*entity.RegisteredProcess, error)
	GlobalSearch(ctx context.Context, filter SearchFilter) ([]*entity.RegisteredProcess, error)
	Update(ctx context.Context, productID string, newRegisteredProcess *entity.RegisteredProcess) error
	GetByID(ctx context.Context, productID string, imageID string) (*entity.RegisteredProcess, error)
}

type SearchFilter struct {
	ProcessType entity.ProcessType
}

func (f SearchFilter) Validate() error {
	return validateFilterProcessType(f.ProcessType)
}

func validateFilterProcessType(processType entity.ProcessType) error {
	if processType != "" {
		return processType.Validate()
	}

	return nil
}
