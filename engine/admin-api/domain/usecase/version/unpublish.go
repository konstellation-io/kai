package version

import (
	"context"
	"fmt"

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
		return nil, fmt.Errorf("error unpublishing version %q: %w", vers.Tag, err)
	}

	vers.PublicationAuthor = nil
	vers.PublicationDate = nil
	vers.Status = entity.VersionStatusStarted

	err = h.versionRepo.Update(productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error updating version %q: %w", vers.Tag, err)
	}

	err = h.userActivityInteractor.RegisterUnpublishAction(user.ID, productID, vers, comment)
	if err != nil {
		return nil, fmt.Errorf("error registering unpublish action: %w", err)
	}

	return vers, nil
}
