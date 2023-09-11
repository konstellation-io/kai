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

var (
	ErrUpdatingVersionStatus = fmt.Errorf("error updating version status")
	ErrUpdatingVersionErrors = fmt.Errorf("error updating version errors")
)

// Start a previously created Version.
func (h *Handler) Start(
	ctx context.Context,
	user *entity.User,
	productID string,
	versionTag string,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	h.logger.Info("Starting version", "userID", user.ID, "versionTag", versionTag, "productID", productID)

	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, nil, err
	}

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, nil, err
	}

	if !vers.CanBeStarted() {
		return nil, nil, internalerrors.ErrInvalidVersionStatusBeforeStarting
	}

	err = h.userActivityInteractor.RegisterStartAction(user.ID, productID, vers, comment)
	if err != nil {
		return nil, nil, err
	}

	versionCfg, err := h.getVersionConfig(ctx, productID, vers)
	if err != nil {
		return nil, nil, err
	}

	vers.Status = entity.VersionStatusStarting

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarting)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarting,
		)
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	go h.startAndNotify(productID, vers, versionCfg, notifyStatusCh)

	return vers, notifyStatusCh, nil
}

func (h *Handler) getVersionConfig(ctx context.Context, productID string, vers *entity.Version) (*entity.VersionConfig, error) {
	versionStreamCfg, err := h.natsManagerService.CreateStreams(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating streams for version %q: %w", vers.Tag, err)
	}

	objectStoreCfg, err := h.natsManagerService.CreateObjectStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating objects stores for version %q: %w", vers.Tag, err)
	}

	kvStoreCfg, err := h.natsManagerService.CreateKeyValueStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating key-value stores for version %q: %w", vers.Tag, err)
	}

	versionCfg := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)

	return versionCfg, nil
}

func (h *Handler) startAndNotify(
	productID string,
	vers *entity.Version,
	versionConfig *entity.VersionConfig,
	notifyStatusCh chan *entity.Version,
) {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
	}()

	err := h.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		h.handleStartError(ctx, productID, vers, notifyStatusCh, err)
		return
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusStarted)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarted,
		)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
}

func (h *Handler) handleStartError(
	ctx context.Context, productID string, vers *entity.Version,
	notifyStatusCh chan *entity.Version, startErr error,
) {
	err := h.versionRepo.SetStatus(ctx, productID, vers.ID, entity.VersionStatusError)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "CRITICAL",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusError,
		)
	}

	errs := []string{startErr.Error()}
	_, err = h.versionRepo.SetErrors(ctx, productID, vers, errs)
	if err != nil {
		h.logger.Error(ErrUpdatingVersionStatus, "ERROR",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusError,
		)
	}

	vers.Status = entity.VersionStatusError
	vers.Errors = append(vers.Errors, startErr.Error())
	notifyStatusCh <- vers
}
