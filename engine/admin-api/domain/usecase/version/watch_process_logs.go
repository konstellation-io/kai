package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (h *Handler) WatchProcessLogs(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag string,
	filters entity.LogFilters,
) (<-chan *entity.ProcessLog, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActViewVersion); err != nil {
		return nil, err
	}

	return h.processLogRepo.WatchProcessLogs(ctx, productID, versionTag, filters)
}
