package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

// Stop removes the resources of the given Version.
func (h *Handler) Stop(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActManageVersion); err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerStopActionFailed(user.Email, productID, v, ErrUserNotAuthorized)

		return nil, nil, err
	}

	h.logger.Info("Stopping version", "userEmail", user.Email, "versionTag", versionTag, "productID", productID)

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerStopActionFailed(user.Email, productID, v, ErrVersionNotFound)

		return nil, nil, err
	}

	if !vers.CanBeStopped() {
		h.registerStopActionFailed(user.Email, productID, vers, ErrVersionCannotBeStopped)
		return nil, nil, ErrVersionCannotBeStopped
	}

	err = h.deleteNatsResources(ctx, productID, vers)
	if err != nil {
		h.registerStopActionFailed(user.Email, productID, vers, ErrDeletingNATSResources)
		return nil, nil, err
	}

	vers.Status = entity.VersionStatusStopping

	err = h.versionRepo.SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStopping)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStopping,
		)
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	go h.stopAndNotify(user.Email, productID, comment, vers, notifyStatusCh)

	return vers, notifyStatusCh, nil
}

func (h *Handler) deleteNatsResources(ctx context.Context, productID string, vers *entity.Version) error {
	err := h.natsManagerService.DeleteStreams(ctx, productID, vers.Tag)
	if err != nil {
		return fmt.Errorf("error deleting stream for version %q: %w", vers.Tag, err)
	}

	err = h.natsManagerService.DeleteObjectStores(ctx, productID, vers.Tag)
	if err != nil {
		return fmt.Errorf("error deleting object stores for version %q: %w", vers.Tag, err)
	}

	err = h.natsManagerService.DeleteVersionKeyValueStores(ctx, productID, vers)
	if err != nil {
		return fmt.Errorf("error deleting key value stores for version %q: %w", vers.Tag, err)
	}

	return nil
}

func (h *Handler) stopAndNotify(
	userEmail,
	productID,
	comment string,
	vers *entity.Version,
	notifyStatusCh chan *entity.Version,
) {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
	}()

	err := h.k8sService.Stop(ctx, productID, vers)
	if err != nil {
		h.registerStopActionFailed(userEmail, productID, vers, ErrStoppingVersion)
		h.handleVersionServiceActionError(ctx, productID, vers, notifyStatusCh, err)

		return
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStopped)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStopped,
		)
	}

	err = h.userActivityInteractor.RegisterStopAction(userEmail, productID, vers, comment)
	if err != nil {
		h.logger.Error(err, "Error registering user activity",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}

	vers.Status = entity.VersionStatusStopped
	notifyStatusCh <- vers
}
