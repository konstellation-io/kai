package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
)

// Unpublish set a Version as not published on DB and K8s.
func (h *Handler) Unpublish(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActUnpublishVersion); err != nil {
		return nil, err
	}

	h.logger.Info(fmt.Sprintf("The user %s is unpublishing version %s on product %s", user.ID, versionTag, productID))

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusPublished {
		return nil, internalerrors.ErrInvalidVersionStatusBeforeUnpublishing
	}

	err = h.k8sService.Unpublish(ctx, productID, v)
	if err != nil {
		return nil, err
	}

	v.PublicationAuthor = nil
	v.PublicationDate = nil
	v.Status = entity.VersionStatusStarted

	err = h.versionRepo.Update(productID, v)
	if err != nil {
		return nil, err
	}

	err = h.userActivityInteractor.RegisterUnpublishAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, err
	}

	return v, nil
}
