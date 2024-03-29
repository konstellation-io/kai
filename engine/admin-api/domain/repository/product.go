package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProductRepo interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	FindAll(ctx context.Context, filter *FindAllFilter) ([]*entity.Product, error)
	FindByIDs(ctx context.Context, productIDs []string, filter *FindAllFilter) ([]*entity.Product, error)
	GetByID(ctx context.Context, productID string) (*entity.Product, error)
	GetByName(ctx context.Context, name string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, productID string) error
	DeleteDatabase(ctx context.Context, productID string) error
}

type FindAllFilter struct {
	ProductName string
}
