package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
)

// SearchByProduct returns all Versions of the given Product, can optionally use a filter.
func (h *Handler) SearchByProduct(
	ctx context.Context,
	user *entity.User,
	productID string,
	filter *repository.ListVersionsFilter,
) ([]*entity.Version, error) {
	err := h.validateFilter(filter)
	if err != nil {
		return nil, err
	}

	versions, err := h.versionRepo.SearchByProduct(ctx, productID, filter)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

func (h *Handler) validateFilter(filter *repository.ListVersionsFilter) error {
	if filter == nil {
		return nil
	}

	if filter.Status != "" {
		err := filter.Status.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
