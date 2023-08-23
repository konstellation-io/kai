package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type MetricRepo interface {
	GetMetrics(
		ctx context.Context,
		startDate time.Time,
		endDate time.Time,
		productID string,
		versionTag string,
	) ([]entity.ClassificationMetric, error)
	CreateIndexes(ctx context.Context, productID string) error
}
