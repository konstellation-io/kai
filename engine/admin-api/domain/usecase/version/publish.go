package version

import (
	"context"
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
) (map[string]string, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActPublishVersion); err != nil {
		return nil, err
	}

	h.logger.Info("Publishing version", "user", user.Email, "product", productID, "version", versionTag)

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, ErrVersionCannotBePublished
	}

	triggerURLs, err := h.k8sService.Publish(ctx, productID, v.Tag)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	v.PublicationDate = &now
	v.PublicationAuthor = &user.Email
	v.Status = entity.VersionStatusPublished

	err = h.versionRepo.Update(productID, v)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"user", user.Email,
			"product", productID,
			"version", versionTag,
		)
	}

	err = h.userActivityInteractor.RegisterPublishAction(user.Email, productID, v, comment)
	if err != nil {
		h.logger.Error(err, "Error registering publish action",
			"user", user.Email,
			"product", productID,
			"version", versionTag,
		)
	}

	return triggerURLs, nil
}
