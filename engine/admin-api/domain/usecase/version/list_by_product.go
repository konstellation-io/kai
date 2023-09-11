package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// ListVersionsByProduct returns all Versions of the given Product.
func (h *Handler) ListVersionsByProduct(ctx context.Context, user *entity.User, productID string) ([]*entity.Version, error) {
	versions, err := h.versionRepo.ListVersionsByProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	return versions, nil
}
