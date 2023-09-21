package version

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
)

// Unpublish set a Version as not published on DB and K8s.
func (h *Handler) Unpublish(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, error) {
	h.logger.Info("Unpublishing version", "userID", user.ID, "versionTag", versionTag, "productID", productID)

	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActUnpublishVersion); err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.ID, productID, v, ErrUserNotAuthorized, "unpublish")

		return nil, err
	}

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.ID, productID, v, ErrVersionNotFound, "unpublish")

		return nil, err
	}

	if vers.Status != entity.VersionStatusPublished {
		h.registerActionFailed(user.ID, productID, vers, ErrVersionCannotBeUnpublished, "unpublish")
		return nil, ErrVersionCannotBeUnpublished
	}

	err = h.k8sService.Unpublish(ctx, productID, vers)
	if err != nil {
		h.registerActionFailed(user.ID, productID, vers, ErrUnpublishingVersion, "unpublish")

		return nil, ErrUnpublishingVersion
	}

	vers.PublicationAuthor = nil
	vers.PublicationDate = nil
	vers.Status = entity.VersionStatusStarted

	err = h.versionRepo.Update(productID, vers)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarted,
		)
	}

	err = h.userActivityInteractor.RegisterUnpublishAction(user.ID, productID, vers, comment)
	if err != nil {
		h.logger.Error(err, "Error registering user activity",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}

	return vers, nil
}
