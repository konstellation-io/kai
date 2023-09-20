package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// GetByVersion returns a Version by its unique tag.
func (h *Handler) GetByVersion(ctx context.Context, user *entity.User, productID, tag string) (*entity.Version, error) {
	return h.versionRepo.GetByVersion(ctx, productID, tag)
}
