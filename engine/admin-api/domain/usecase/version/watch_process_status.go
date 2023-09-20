package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

func (h *Handler) WatchProcessStatus(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag string,
) (<-chan *entity.Process, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct); err != nil {
		return nil, err
	}

	v, err := h.versionRepo.GetByVersion(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	return h.k8sService.WatchProcessStatus(ctx, productID, v.Tag)
}
