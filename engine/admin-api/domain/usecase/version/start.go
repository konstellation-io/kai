package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

// Start a previously created Version.
func (h *Handler) Start(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, chan *entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.Email, productID, v, ErrUserNotAuthorized, StartAction)

		return nil, nil, err
	}

	h.logger.Info("Starting version", "userEmail", user.Email, "versionTag", versionTag, "productID", productID)

	vers, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		v := &entity.Version{Tag: versionTag}
		h.registerActionFailed(user.Email, productID, v, ErrVersionNotFound, StartAction)

		return nil, nil, err
	}

	if !vers.CanBeStarted() {
		h.registerActionFailed(user.Email, productID, vers, ErrVersionCannotBeStarted, StartAction)
		return nil, nil, ErrVersionCannotBeStarted
	}

	versionCfg, err := h.createStreamingResources(ctx, productID, vers)
	if err != nil {
		h.registerActionFailed(user.Email, productID, vers, ErrCreatingNATSResources, StartAction)
		return nil, nil, err
	}

	err = h.updateKeyValueConfigurations(ctx, vers, versionCfg)
	if err != nil {
		return nil, nil, err
	}

	vers.Status = entity.VersionStatusStarting

	err = h.versionRepo.SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarting)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarting,
		)
	}

	notifyStatusCh := make(chan *entity.Version, 1)

	go h.startAndNotify(user.Email, productID, comment, vers, versionCfg, notifyStatusCh)

	return vers, notifyStatusCh, nil
}

func (h *Handler) updateKeyValueConfigurations(
	ctx context.Context,
	vers *entity.Version,
	versionCfg *entity.VersionStreamingResources,
) error {
	// Version kv store
	var kvConfigurations []entity.KeyValueConfiguration

	if len(vers.Config) > 0 {
		kvConfigurations = append(kvConfigurations, entity.KeyValueConfiguration{
			Store:         versionCfg.KeyValueStores.KeyValueStore,
			Configuration: vers.Config,
		})
	}

	// Workflows configuration
	for _, workflow := range vers.Workflows {
		workflowKVstore, err := versionCfg.KeyValueStores.GetWorkflowKeyValueStore(workflow.Name)
		if err != nil {
			return err
		}

		if len(workflow.Config) > 0 {
			kvConfigurations = append(kvConfigurations, entity.KeyValueConfiguration{
				Store:         workflowKVstore,
				Configuration: workflow.Config,
			})
		}

		// Processes configuration
		for _, process := range workflow.Processes {
			processKVStore, err := versionCfg.KeyValueStores.Workflows[workflow.Name].GetProcessKeyValueStore(process.Name)
			if err != nil {
				return err
			}

			if len(process.Config) > 0 {
				kvConfigurations = append(kvConfigurations, entity.KeyValueConfiguration{
					Store:         processKVStore,
					Configuration: process.Config,
				})
			}
		}
	}

	if len(kvConfigurations) > 0 {
		err := h.natsManagerService.UpdateKeyValueConfiguration(ctx, kvConfigurations)
		if err != nil {
			return fmt.Errorf("updating key-value configurations: %w", err)
		}
	}

	return nil
}

func (h *Handler) createStreamingResources(
	ctx context.Context,
	productID string,
	vers *entity.Version,
) (*entity.VersionStreamingResources, error) {
	versionStreamCfg, err := h.natsManagerService.CreateStreams(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating streams for version %q: %w", vers.Tag, err)
	}

	objectStoreCfg, err := h.natsManagerService.CreateObjectStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating objects stores for version %q: %w", vers.Tag, err)
	}

	kvStoreCfg, err := h.natsManagerService.CreateVersionKeyValueStores(ctx, productID, vers)
	if err != nil {
		return nil, fmt.Errorf("error creating key-value stores for version %q: %w", vers.Tag, err)
	}

	return entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)
}

func (h *Handler) startAndNotify(
	userEmail,
	productID,
	comment string,
	vers *entity.Version,
	versionConfig *entity.VersionStreamingResources,
	notifyStatusCh chan *entity.Version,
) {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer func() {
		cancel()
		close(notifyStatusCh)
	}()

	err := h.k8sService.Start(ctx, productID, vers, versionConfig)
	if err != nil {
		h.registerActionFailed(userEmail, productID, vers, ErrStartingVersion, StartAction)
		h.handleVersionServiceActionError(ctx, productID, vers, notifyStatusCh, err)

		return
	}

	err = h.versionRepo.SetStatus(ctx, productID, vers.Tag, entity.VersionStatusStarted)
	if err != nil {
		h.logger.Error(err, "Error updating version status",
			"productID", productID,
			"versionTag", vers.Tag,
			"previousStatus", vers.Status,
			"newStatus", entity.VersionStatusStarted,
		)
	}

	err = h.userActivityInteractor.RegisterStartAction(userEmail, productID, vers, comment)
	if err != nil {
		h.logger.Error(err, "Error registering user activity",
			"productID", productID,
			"versionTag", vers.Tag,
			"comment", comment,
		)
	}

	vers.Status = entity.VersionStatusStarted
	notifyStatusCh <- vers
}
