package logs

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
)

var _ LogsUsecase = (*LogsInteractor)(nil)

type LogsUsecase interface {
	GetLogs(logFilters entity.LogFilters) ([]*entity.Log, error)
}

type LogsInteractor struct {
	logsService service.LogsService
}

func NewLogsInteractor(logsService service.LogsService) *LogsInteractor {
	return &LogsInteractor{
		logsService,
	}
}

func (i *LogsInteractor) GetLogs(logFilters entity.LogFilters) ([]*entity.Log, error) {
	return i.logsService.GetLogs(logFilters)
}
