package repository

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type ProcessLogRepository interface {
	WatchProcessLogs(ctx context.Context, runtimeID, versionName string, filters entity.LogFilters) (<-chan *entity.ProcessLog, error)
	PaginatedSearch(ctx context.Context, runtimeID string, searchOpts *entity.SearchLogsOptions) (*entity.SearchLogsResult, error)
	CreateIndexes(ctx context.Context, runtimeID string) error
}
