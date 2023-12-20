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
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActUnpublishVersion); err != nil {
		return nil, err
	}

	h.logger.Info("Unpublishing version", "userEmail", user.Email, "versionTag", versionTag, "productID", productID)

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if vers.Status != entity.VersionStatusPublished {
		return nil, ErrVersionCannotBeUnpublished
	}

	err = h.k8sService.Unpublish(ctx, productID, vers)
	if err != nil {
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

	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	product.RemovePublishedVersion()

	err = h.productRepo.Update(context.Background(), product)
	if err != nil {
		h.logger.Error(err, "Error updating product published version",
			"productID", productID,
			"publishedVersion", vers.Tag,
		)
	}

	err = h.userActivityInteractor.RegisterUnpublishAction(user.Email, productID, vers, comment)
	if err != nil {
		h.logger.Error(err, "Error registering user activity",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}

	return vers, nil
}
