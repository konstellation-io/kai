package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// GetByTag returns a Version by its unique tag.
func (h *Handler) GetByTag(ctx context.Context, user *entity.User, productID, tag string) (*entity.Version, error) {
	return h.versionRepo.GetByTag(ctx, productID, tag)
}
