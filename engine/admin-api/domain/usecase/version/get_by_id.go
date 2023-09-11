package version

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

// GetByID returns a Version by its unique ID.
func (h *Handler) GetByID(productID, versionID string) (*entity.Version, error) {
	return h.versionRepo.GetByID(productID, versionID)
}
