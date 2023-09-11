package version

import (
	"context"
	"fmt"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (h *Handler) SearchLogs(
	ctx context.Context,
	user *entity.User,
	productID string,
	filters entity.LogFilters,
	cursor *string,
) (*entity.SearchLogsResult, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActViewVersion); err != nil {
		return nil, err
	}

	startDate, err := time.Parse(time.RFC3339, filters.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	var endDate time.Time
	if filters.EndDate != nil {
		endDate, err = time.Parse(time.RFC3339, *filters.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
	} else {
		endDate = time.Now()
	}

	options := &entity.SearchLogsOptions{
		Cursor:         cursor,
		StartDate:      startDate,
		EndDate:        endDate,
		Search:         filters.Search,
		ProcessIDs:     filters.ProcessIDs,
		Levels:         filters.Levels,
		VersionsIDs:    filters.VersionsIDs,
		WorkflowsNames: filters.WorkflowsNames,
	}

	return h.processLogRepo.PaginatedSearch(ctx, productID, options)
}
