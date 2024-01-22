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
	ErrVersionAlreadyPublished = errors.New("version already published")
	ErrVersionBeingPublished   = errors.New("version is already being published")
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

	v, err := h.versionRepo.GetByTag(ctx, params.ProductID, params.VersionTag)
	if err != nil {
		return nil, nil, err
	}

	if v.Status == entity.VersionStatusPublished {
		return nil, nil, ErrVersionAlreadyPublished
	}

	if v.Status == entity.VersionStatusPublishing {
		return nil, nil, ErrVersionBeingPublished
	}

	if !params.Force {
		if product.HasVersionPublished() {
			return nil, nil, ErrProductAlreadyPublished
		}

		if v.Status != entity.VersionStatusStarted {
			return nil, nil, ErrVersionIsNotStarted
		}
	}

	compensations := compensator.New()

	notifyCh := make(chan *entity.Version, 1)

	go func() {
		defer close(notifyCh)

		err := h.publishVersion(compensations, user, product, v, params.Comment, params.Force)
		if err != nil {
			h.handleAsyncVersionError(compensations, product.ID, v, err)
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
	force bool,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer cancel()

	if force {
		err := h.forcePublishPreviousActions(ctx, compensations, user, product, version, comment)
		if err != nil {
			return err
		}
	}

	err := h.k8sService.Publish(ctx, product.ID, version.Tag)
	if err != nil {
		return err
	}

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

func (h *Handler) forcePublishPreviousActions(
	ctx context.Context,
	compensations *compensator.Compensator,
	user *entity.User,
	product *entity.Product,
	version *entity.Version,
	comment string,
) error {
	if product.HasVersionPublished() {
		err := h.tearDownPreviousVersion(ctx, compensations, user, product, comment)
		if err != nil {
			return err
		}
	}

	if version.Status != entity.VersionStatusStarted {
		err := h.startVersionSync(ctx, user, product.ID, version.Tag, comment)
		if err != nil {
			return err
		}

		compensations.AddCompensation(func() error { _, _, err := h.Stop(ctx, user, product.ID, version.Tag, comment); return err })
	}

	return nil
}

func (h *Handler) tearDownPreviousVersion(
	ctx context.Context,
	compensations *compensator.Compensator,
	user *entity.User,
	product *entity.Product,
	comment string,
) error {
	publishedVersion := *product.PublishedVersion

	_, err := h.Unpublish(ctx, user, product.ID, publishedVersion, comment)
	if err != nil {
		return fmt.Errorf("unpublishing previous version: %w", err)
	}

	compensations.AddCompensation(func() error {
		_, _, err := h.Publish(ctx, user, PublishParams{
			ProductID:  product.ID,
			VersionTag: publishedVersion,
			Comment:    comment,
			Force:      true,
		})
		return err
	})

	// _, _, err = h.Stop(ctx, user, product.ID, publishedVersion, comment)
	// if err != nil {
	// 	 return fmt.Errorf("stopping previous version: %w", err)
	// }

	compensations.AddCompensation(h.startVersionSyncFunc(user, product.ID, publishedVersion, comment))

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

func (h *Handler) startVersionSync(ctx context.Context, user *entity.User, productID, versionTag, comment string) error {
	_, startNotifyCh, err := h.Start(ctx, user, productID, versionTag, comment)
	if err != nil {
		return fmt.Errorf("start version: %w", err)
	}

	startedVersion := <-startNotifyCh

	if startedVersion.Status == entity.VersionStatusError || startedVersion.Status == entity.VersionStatusCritical {
		return fmt.Errorf("starting version: %s", startedVersion.Error) //nolint:goerr113 // future refactor
	}

	return nil
}

func (h *Handler) startVersionSyncFunc(user *entity.User, productID, publishedVersion, comment string) compensator.Compensation {
	return func() error {
		return h.startVersionSync(context.Background(), user, productID, publishedVersion, comment)
	}
}
