package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

// GetByTag returns a Version by its unique tag.
func (h *Handler) GetByTag(ctx context.Context, user *entity.User, productID, tag string) (*entity.Version, error) {
	_, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	err = h.accessControl.CheckProductGrants(user, productID, auth.ActViewProduct)
	if err != nil {
		return nil, err
	}

	version, err := h.versionRepo.GetByTag(ctx, productID, tag)
	if err != nil {
		return nil, err
	}

	if version.Status != entity.VersionStatusPublished {
		return version, err
	}

	publishedTriggers, err := h.k8sService.GetPublishedTriggers(ctx, productID)
	if err != nil {
		return nil, err
	}

	version.PublishedTriggers = publishedTriggers

	return version, nil
}
