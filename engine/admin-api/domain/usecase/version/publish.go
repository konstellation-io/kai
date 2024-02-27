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
	ErrVersionAlreadyPublished = errors.New("version already published")
)

type PublishOpts struct {
	ProductID  string
	VersionTag string
	Comment    string
	Force      bool
}

// Publish set a Version as published on DB and K8s.
func (h *Handler) Publish(ctx context.Context, user *entity.User, opts PublishOpts) (map[string]string, error) {
	if err := h.accessControl.CheckProductGrants(user, opts.ProductID, auth.ActManageVersion); err != nil {
		return nil, err
	}

	h.logger.Info("Publishing version", "user", user.Email, "product", opts.ProductID, "version", opts.VersionTag)

	product, err := h.productRepo.GetByID(ctx, opts.ProductID)
	if err != nil {
		return nil, err
	}

	version, err := h.versionRepo.GetByTag(ctx, opts.ProductID, opts.VersionTag)
	if err != nil {
		return nil, err
	}

	if version.Status == entity.VersionStatusPublished {
		return nil, ErrVersionAlreadyPublished
	}

	if version.Status != entity.VersionStatusStarted {
		return nil, ErrVersionIsNotStarted
	}

	compensations := compensator.New()

	if product.HasVersionPublished() {
		if !opts.Force {
			return nil, ErrProductAlreadyPublished
		}

		publishedVersion := *product.PublishedVersion

		err := h.versionRepo.SetStatus(ctx, product.ID, publishedVersion, entity.VersionStatusStarted)
		if err != nil {
			return nil, err
		}

		compensations.AddCompensation(func() error {
			return h.versionRepo.SetStatus(context.Background(), product.ID, publishedVersion, entity.VersionStatusPublished)
		})
	}

	urls, err := h.publishVersion(ctx, compensations, user, product, version, opts.Comment)
	if err != nil {
		go h.handlePublicationError(err, compensations, product.ID, version)
		return nil, err
	}

	return urls, nil
}

func (h *Handler) publishVersion(
	ctx context.Context,
	compensations *compensator.Compensator,
	user *entity.User,
	product *entity.Product,
	version *entity.Version,
	comment string,
) (map[string]string, error) {
	triggerURLs, err := h.k8sService.Publish(ctx, product.ID, version.Tag)
	if err != nil {
		return nil, err
	}

	compensations.AddCompensation(h.rollbackPublishedVersionFunc(product, version))

	err = h.updateVersionStatusToPublished(compensations, user, product, version)
	if err != nil {
		return nil, err
	}

	err = h.updateProductPublishedVersion(ctx, compensations, product, version)
	if err != nil {
		return nil, err
	}

	err = h.userActivityInteractor.RegisterPublishAction(user.Email, product.ID, version, comment)
	if err != nil {
		return nil, fmt.Errorf("registering publish action: %w", err)
	}

	return triggerURLs, nil
}

func (h *Handler) updateVersionStatusToPublished(
	compensations *compensator.Compensator,
	user *entity.User,
	product *entity.Product,
	version *entity.Version,
) error {
	version.SetPublishStatus(user.Email)

	err := h.versionRepo.Update(product.ID, version)
	if err != nil {
		return fmt.Errorf("updating version: %w", err)
	}

	compensations.AddCompensation(h.setVersionUnpublishStatusFunc(product.ID, version))

	return nil
}

func (h *Handler) updateProductPublishedVersion(
	ctx context.Context,
	compensations *compensator.Compensator,
	product *entity.Product,
	version *entity.Version,
) error {
	rollbackProductPublishedVersionFunc := h.rollbackProductPublishedVersionFunc(product)

	product.UpdatePublishedVersion(version.Tag)

	err := h.productRepo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("updating product's published version: %w", err)
	}

	compensations.AddCompensation(rollbackProductPublishedVersionFunc)

	return nil
}

func (h *Handler) rollbackPublishedVersionFunc(product *entity.Product, version *entity.Version) compensator.Compensation {
	if product.HasVersionPublished() {
		previouslyPublishedVersion := *product.PublishedVersion

		return func() error {
			_, err := h.k8sService.Publish(context.Background(), product.ID, previouslyPublishedVersion)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return func() error {
		return h.k8sService.Unpublish(context.Background(), product.ID, version)
	}
}

func (h *Handler) setVersionUnpublishStatusFunc(productID string, version *entity.Version) compensator.Compensation {
	return func() error {
		version.UnsetPublishStatus()

		return h.versionRepo.Update(productID, version)
	}
}

func (h *Handler) rollbackProductPublishedVersionFunc(product *entity.Product) compensator.Compensation {
	if product.HasVersionPublished() {
		previouslyPublishedVersion := *product.PublishedVersion

		return func() error {
			product.UpdatePublishedVersion(previouslyPublishedVersion)

			return h.productRepo.Update(context.Background(), product)
		}
	}

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
