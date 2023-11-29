package service

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/service_${GOFILE} -package=mocks

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

type LogsService interface {
	GetLogs(logFilters entity.LogFilters) ([]*entity.Log, error)
}
