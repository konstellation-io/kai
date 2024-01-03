package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// GetLatest returns a the latest created version of a product.
func (h *Handler) GetLatest(ctx context.Context, user *entity.User, productID string) (*entity.Version, error) {
	return h.versionRepo.GetLatest(ctx, productID)
}
