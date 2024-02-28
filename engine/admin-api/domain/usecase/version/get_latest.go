package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// GetLatest returns a the latest created version of a product.
func (h *Handler) GetLatest(ctx context.Context, user *entity.User, productID string) (*entity.Version, error) {
	_, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	ver, err := h.versionRepo.GetLatest(ctx, productID)
	if err != nil {
		return nil, ErrProductExistsNoVersions
	}

	return ver, nil
}
