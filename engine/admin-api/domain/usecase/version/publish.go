package version

import (
	"context"
	"fmt"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

// Publish set a Version as published on DB and K8s.
func (h *Handler) Publish(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, error) {
	h.logger.Info("Publishing version", "userID", user.ID, "versionTag", versionTag, "productID", productID)

	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActPublishVersion); err != nil {
		return nil, err
	}

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, ErrVersionCannotBePublished
	}

	err = h.k8sService.Publish(ctx, productID, v)
	if err != nil {
		return nil, err
	}

	previousPublishedVersion, err := h.versionRepo.ClearPublishedVersion(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error unpublishing previous version: %w", err)
	}

	now := time.Now()
	v.PublicationDate = &now
	v.PublicationAuthor = &user.ID
	v.Status = entity.VersionStatusPublished

	err = h.versionRepo.Update(productID, v)
	if err != nil {
		return nil, err
	}

	err = h.userActivityInteractor.RegisterPublishAction(user.ID, productID, v, previousPublishedVersion, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}
