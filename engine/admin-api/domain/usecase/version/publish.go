package version

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/pkg/compensator"
	"github.com/spf13/viper"
)

var (
	ErrProductAlreadyPublished = errors.New("product already has a published version")
)

type PublishParams struct {
	ProductID  string
	VersionTag string
	Comment    string
	Force      bool
}

// Publish set a Version as published on DB and K8s.
func (h *Handler) Publish(
	ctx context.Context,
	user *entity.User,
	params PublishParams,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, params.ProductID, auth.ActPublishVersion); err != nil {
		return nil, nil, err
	}

	h.logger.Info("Publishing version", "user", user.Email, "product", params.ProductID, "version", params.VersionTag)

	product, err := h.productRepo.GetByID(ctx, params.ProductID)
	if err != nil {
		return nil, nil, err
	}

	if product.HasVersionPublished() {
		return nil, nil, ErrProductAlreadyPublished
	}

	v, err := h.versionRepo.GetByTag(ctx, params.ProductID, params.VersionTag)
	if err != nil {
		return nil, nil, err
	}

	if v.Status != entity.VersionStatusStarted {
		return nil, nil, ErrVersionCannotBePublished
	}

	compensations := compensator.New()

	v.SetPublishingStatus(user.Email)

	compensations.AddCompensation(func() error { v.SetStartedStatus(); return nil })

	notifyCh := make(chan *entity.Version, 1)

	go func() {
		defer close(notifyCh)

		err := h.publishVersion(compensations, user, product, v, params.Comment)
		if err != nil {
			h.handlePublicationError(err, compensations, product.ID, v)
		}

		notifyCh <- v
	}()

	return v, notifyCh, nil
}

func (h *Handler) publishVersion(
	compensations *compensator.Compensator,
	user *entity.User,
	product *entity.Product,
	version *entity.Version,
	comment string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer cancel()

	triggerURLs, err := h.k8sService.Publish(ctx, product.ID, version.Tag)
	if err != nil {
		return err
	}

	_ = triggerURLs

	compensations.AddCompensation(h.unpublishVersionFunc(product.ID, version))

	version.SetPublishStatus(user.Email)

	err = h.versionRepo.Update(product.ID, version)
	if err != nil {
		return fmt.Errorf("updating version: %w", err)
	}

	compensations.AddCompensation(h.setVersionUnpublishStatusFunc(product.ID, version))

	product.UpdatePublishedVersion(version.Tag)

	err = h.productRepo.Update(ctx, product)
	if err != nil {
		return fmt.Errorf("updating product's published version: %w", err)
	}

	compensations.AddCompensation(h.removeProductPublishedVersionFunc(product))

	err = h.userActivityInteractor.RegisterPublishAction(user.Email, product.ID, version, comment)
	if err != nil {
		return fmt.Errorf("registering publish action: %w", err)
	}

	return nil
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
		"product", productID,
		"version", version.Tag,
	)

	err := compensations.Execute()
	if err != nil {
		h.handleCriticalError(context.Background(), productID, version, err)
	}
}
