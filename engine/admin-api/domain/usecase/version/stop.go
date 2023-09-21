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
	h.logger.Info("Stopping version", "userID", user.ID, "versionTag", versionTag, "productID", productID)

	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStopVersion); err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.ID, productID, v, ErrUserNotAuthorized, "stop")

		return nil, nil, err
	}

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.ID, productID, v, ErrVersionNotFound, "stop")

		return nil, nil, err
	}

	if !vers.CanBeStopped() {
		h.registerActionFailed(user.ID, productID, vers, ErrVersionCannotBeStopped, "stop")
		return nil, nil, ErrVersionCannotBeStopped
	}

	err = h.deleteNatsResources(ctx, productID, vers)
	if err != nil {
		h.registerActionFailed(user.ID, productID, vers, ErrDeletingNATSResources, "stop")
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

	go h.stopAndNotify(user.ID, productID, comment, vers, notifyStatusCh)

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

	return nil
}

func (h *Handler) stopAndNotify(
	userID,
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
		h.registerActionFailed(userID, productID, vers, ErrStoppingVersion, "stop")
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

	err = h.userActivityInteractor.RegisterStopAction(userID, productID, vers, comment)
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
