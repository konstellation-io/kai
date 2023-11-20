package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/pkg/compensator"
)

var (
	ErrProductAlreadyPublished = errors.New("product already has a published version")
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

	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if product.HasVersionPublished() {
		return nil, ErrProductAlreadyPublished
	}

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, ErrVersionCannotBePublished
	}

	compensations := compensator.New()

	triggerURLs, err := h.k8sService.Publish(ctx, productID, v.Tag)
	if err != nil {
		return nil, err
	}

	compensations.AddCompensation(h.unpublishVersionFunc(productID, v))

	v.SetPublishStatus(user.Email)

	err = h.versionRepo.Update(productID, v)
	if err != nil {
		go h.handlePublicationError(err, compensations, productID, v)
		return nil, fmt.Errorf("updating version: %w", err)
	}

	compensations.AddCompensation(h.setVersionUnpublishStatusFunc(productID, v))

	product.UpdatePublishedVersion(v.Tag)

	err = h.productRepo.Update(ctx, product)
	if err != nil {
		go h.handlePublicationError(err, compensations, productID, v)
		return nil, fmt.Errorf("updating product's published version: %w", err)
	}

	compensations.AddCompensation(h.removeProductPublishedVersionFunc(product))

	err = h.userActivityInteractor.RegisterPublishAction(user.Email, productID, v, comment)
	if err != nil {
		go h.handlePublicationError(err, compensations, productID, v)
		return nil, fmt.Errorf("registering publish action: %w", err)
	}

	return triggerURLs, nil
}

func (h *Handler) unpublishVersionFunc(productID string, version *entity.Version) compensator.Compensation {
	return func() error {
		return h.k8sService.Unpublish(context.Background(), productID, version)
	}
}

func (h *Handler) setVersionUnpublishStatusFunc(productID string, version *entity.Version) compensator.Compensation {
	return func() error {
		version.UnsetPublishStatus()

		return h.versionRepo.Update(productID, version)
	}
}

func (h *Handler) removeProductPublishedVersionFunc(product *entity.Product) compensator.Compensation {
	return func() error {
		product.RemovePublishedVersion()

		return h.productRepo.Update(context.Background(), product)
	}
}

func (h *Handler) handlePublicationError(
	publicationError error,
	compensations *compensator.Compensator,
	productID string,
	version *entity.Version,
) {
	h.logger.Error(publicationError, "Error during version publication, executing compensations...",
		"productID", productID,
		"versionTag", version.Tag,
	)

	err := compensations.Execute()
	if err != nil {
		h.handleCriticalError(context.Background(), productID, version, err)
	}
}
