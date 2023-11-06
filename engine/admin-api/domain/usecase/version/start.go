package version

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/pkg/compensator"
	"github.com/spf13/viper"
)

// Start a previously created Version.
func (h *Handler) Start(
	ctx context.Context,
	user *entity.User,
	productID,
	versionTag,
	comment string,
) (*entity.Version, error) {
	if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartVersion); err != nil {
		return nil, err
	}

	h.logger.Info("Starting version", "userEmail", user.Email, "versionTag", versionTag, "productID", productID)

	compensations := compensator.New()

	product, err := h.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	version, err := h.versionRepo.GetByTag(ctx, productID, versionTag)
	if err != nil {
		return nil, err
	}

	if version.Status == entity.VersionStatusCritical {
		if err := h.accessControl.CheckProductGrants(user, productID, auth.ActStartCriticalVersion); err != nil {
			return nil, err
		}
	}

	if !version.CanBeStarted() {
		return nil, ErrVersionCannotBeStarted
	}

	version.Status = entity.VersionStatusStarting

	err = h.versionRepo.SetStatus(ctx, productID, version.Tag, entity.VersionStatusStarting)
	if err != nil {
		return nil, fmt.Errorf("setting version status to %q: %w", entity.VersionStatusStarting, err)
	}

	go func() {
		err = h.createVersionResources(user, product, version, comment, compensations)
		if err != nil {
			h.handleStartVersionError(productID, version, err, compensations)
			return
		}
	}()

	return version, nil
}

func (h *Handler) createVersionResources(
	user *entity.User,
	product *entity.Product,
	version *entity.Version,
	comment string,
	compensations *compensator.Compensator,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration(config.VersionStatusTimeoutKey))
	defer cancel()

	versionStreamCfg, err := h.natsManagerService.CreateStreams(ctx, product.ID, version)
	if err != nil {
		return fmt.Errorf("error creating streams for version %q: %w", version.Tag, err)
	}

	compensations.AddCompensation(h.deleteStreamFunc(product.ID, version))

	objectStoreCfg, err := h.natsManagerService.CreateObjectStores(ctx, product.ID, version)
	if err != nil {
		return fmt.Errorf("error creating objects stores for version %q: %w", version.Tag, err)
	}

	compensations.AddCompensation(h.deleteObjectStoresFunc(product.ID, version))

	kvStoreCfg, err := h.natsManagerService.CreateVersionKeyValueStores(ctx, product.ID, version)
	if err != nil {
		return fmt.Errorf("error creating key-value stores for version %q: %w", version.Tag, err)
	}

	kvStoreCfg.GlobalKeyValueStore = product.KeyValueStore

	versionCfg, err := entity.NewVersionConfig(versionStreamCfg, objectStoreCfg, kvStoreCfg)
	if err != nil {
		return err
	}

	err = h.updateKeyValueConfigurations(ctx, version, versionCfg)
	if err != nil {
		return fmt.Errorf("initializing centralized configuration: %w", err)
	}

	err = h.k8sService.Start(ctx, product, version, versionCfg)
	if err != nil {
		return fmt.Errorf("starting version on k8s service: %w", err)
	}

	compensations.AddCompensation(h.stopVersionFunc(product.ID, version))

	err = h.versionRepo.SetStatus(ctx, product.ID, version.Tag, entity.VersionStatusStarted)
	if err != nil {
		return fmt.Errorf("updating version status to %q: %w", entity.VersionStatusStarted, err)
	}

	err = h.userActivityInteractor.RegisterStartAction(user.Email, product.ID, version, comment)
	if err != nil {
		return fmt.Errorf("registering start action: %w", err)
	}

	return nil
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
			Store:         versionCfg.KeyValueStores.VersionKeyValueStore,
			Configuration: vers.Config,
		})
	}

	// Workflows configuration
	for _, workflow := range vers.Workflows {
		workflowConfigurations, err := h.getWorkflowConfigurations(versionCfg, workflow)
		if err != nil {
			return err
		}

		kvConfigurations = append(kvConfigurations, workflowConfigurations...)
	}

	if len(kvConfigurations) > 0 {
		err := h.natsManagerService.UpdateKeyValueConfiguration(ctx, kvConfigurations)
		if err != nil {
			return fmt.Errorf("updating key-value configurations: %w", err)
		}
	}

	return nil
}

func (h *Handler) handleStartVersionError(
	productID string,
	version *entity.Version,
	startError error,
	compensations *compensator.Compensator,
) {
	h.logger.Error(startError, "Error starting version", "productID", productID, "versionTag", version.Tag)

	ctx := context.Background()

	err := compensations.Execute()
	if err != nil {
		h.handleCriticalError(ctx, productID, version, err)
		return
	}

	err = h.versionRepo.SetErrorStatusWithError(ctx, productID, version.Tag, startError.Error())
	if err != nil {
		h.logger.Error(err, "Updating version with error", "productID", productID, "versionTag", version.Tag)
	}
}

func (h *Handler) handleCriticalError(ctx context.Context, productID string, version *entity.Version, err error) {
	err = h.versionRepo.SetCriticalStatusWithError(ctx, productID, version.Tag, err.Error())
	if err != nil {
		h.logger.Error(err,
			"Error setting status version",
			"productID", productID, "versionTag", version.Tag, "wantedStatus", entity.VersionStatusCritical,
		)
	}
}

func (h *Handler) getWorkflowConfigurations(
	versionCfg *entity.VersionStreamingResources,
	workflow entity.Workflow,
) ([]entity.KeyValueConfiguration, error) {
	var workflowConfigurations []entity.KeyValueConfiguration

	workflowKVstore, err := versionCfg.KeyValueStores.GetWorkflowKeyValueStore(workflow.Name)
	if err != nil {
		return nil, err
	}

	if len(workflow.Config) > 0 {
		workflowConfigurations = append(workflowConfigurations, entity.KeyValueConfiguration{
			Store:         workflowKVstore,
			Configuration: workflow.Config,
		})
	}

	// Processes configuration
	for _, process := range workflow.Processes {
		processKVStore, err := versionCfg.KeyValueStores.Workflows[workflow.Name].GetProcessKeyValueStore(process.Name)
		if err != nil {
			return nil, err
		}

		if len(process.Config) > 0 {
			workflowConfigurations = append(workflowConfigurations, entity.KeyValueConfiguration{
				Store:         processKVStore,
				Configuration: process.Config,
			})
		}
	}

	return workflowConfigurations, nil
}

func (h *Handler) deleteStreamFunc(productID string, version *entity.Version) func() error {
	return func() error {
		return h.natsManagerService.DeleteStreams(context.Background(), productID, version.Tag)
	}
}

func (h *Handler) deleteObjectStoresFunc(productID string, version *entity.Version) func() error {
	return func() error {
		return h.natsManagerService.DeleteObjectStores(context.Background(), productID, version.Tag)
	}
}

func (h *Handler) stopVersionFunc(productID string, version *entity.Version) func() error {
	return func() error {
		return h.k8sService.Stop(context.Background(), productID, version)
	}
}
