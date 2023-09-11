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

// Start create the resources of the given Version.
func (h *Handler) Start(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionTag string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, nil, err
	}

	h.logger.Info(fmt.Sprintf("The user %q is starting version %q on product %q", user.ID, versionTag, productID))

	v, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, nil, err
	}

	if !v.CanBeStarted() {
		return nil, nil, internalerrors.ErrInvalidVersionStatusBeforeStarting
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	err = h.versionRepo.SetStatus(ctx, productID, v.ID, entity.VersionStatusStarting)
	if err != nil {
		return nil, nil, err
	}

	// Notify intermediate state
	v.Status = entity.VersionStatusStarting
	notifyStatusCh <- v

	err = h.userActivityInteractor.RegisterStartAction(user.ID, productID, v, comment)
	if err != nil {
		return nil, nil, err
	}

	versionStreamCfg, err := h.natsManagerService.CreateStreams(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating streams for version %q: %w", v.Tag, err)
	}

	objectStoreCfg, err := h.natsManagerService.CreateObjectStores(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating objects stores for version %q: %w", v.Tag, err)
	}

	kvStoreCfg, err := h.natsManagerService.CreateKeyValueStores(ctx, productID, v)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating key-value stores for version %q: %w", v.Tag, err)
	}

	versionCfg := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)

	go h.startAndNotify(productID, v, versionCfg, notifyStatusCh)

	return v, notifyStatusCh, nil
}

func (h *Handler) startAndNotify(
	productID string,
	vers *entity.Version,
	versionConfig *entity.VersionConfig,
	notifyStatusCh chan *entity.Version,
) {
	// WARNING: This function doesn't handle error because there is no  ERROR status defined for a Version
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
		h.logger.V(0).Info("[versionInteractor.startAndNotify] channel closed")
	}()

	err := h.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		h.logger.Error(err, "[versionInteractor.startAndNotify] error starting version", "version tag", vers.Tag)
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarted)
	if err != nil {
		h.logger.Error(err, "[versionInteractor.startAndNotify] error starting version", "version tag", vers.Tag)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
	h.logger.Info("[versionInteractor.startAndNotify] version started", "version tag", vers.Tag)
}
