package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	internalerrors "github.com/konstellation-io/kai/engine/admin-api/internal/errors"
	"github.com/spf13/viper"
)

// Stop removes the resources of the given Version.
func (h *Handler) Stop(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionTag string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, nil, err
	}

	h.logger.Info(fmt.Sprintf("The user %q is stopping version %q on product %q", user.ID, versionTag, productID))

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStopped() {
		return nil, nil, internalerrors.ErrInvalidVersionStatusBeforeStopping
	}

	err = h.versionRepo.SetStatus(ctx, productID, v.ID, entity.VersionStatusStopping)
	if err != nil {
		return nil, nil, err
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	// Notify intermediate state
	v.Status = entity.VersionStatusStopping
	notifyStatusCh <- v

	err = h.userActivityInteractor.RegisterStopAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	err = h.natsManagerService.DeleteStreams(ctx, productID, versionTag)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionTag, err)
	}

	err = h.natsManagerService.DeleteObjectStores(ctx, productID, versionTag)
	if err != nil {
		return nil, nil, fmt.Errorf("error stopping version %q: %w", versionTag, err)
	}

	go h.stopAndNotify(productID, v, notifyStatusCh)

	return v, notifyStatusCh, nil
}

func (h *Handler) stopAndNotify(
	productID string,
	vers *entity.Version,
	notifyStatusCh chan *entity.Version,
) {
	// WARNING: This function doesn't handle error because there is no  ERROR status defined for a Version
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
		h.logger.V(0).Info("[versionInteractor.stopAndNotify] channel closed")
	}()

	err := h.k8sService.Stop(ctx, productID, vers)
	if err != nil {
		h.logger.Error(err, "[versionInteractor.stopAndNotify] error stopping version", "version tag", vers.Tag)
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStopped)
	if err != nil {
		h.logger.Error(err, "[versionInteractor.stopAndNotify] error stopping version", "version tag", vers.Tag)
	}

	vers.Status = entity.VersionStatusStopped
	notifyStatusCh <- vers
	h.logger.Info("[versionInteractor.stopAndNotify] version stopped", "version tag", vers.Tag)
}
