package process

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/spf13/viper"
)

func (ps *Service) RegisterProcess(
	ctx context.Context,
	user *entity.User,
	opts RegisterProcessOpts,
) (*entity.RegisteredProcess, error) {
	ps.logger.Info("Registering new process")

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if err := ps.checkGrants(user, opts); err != nil {
		return nil, err
	}

	scope := ps.getProcessRegisterScope(opts)
	processToRegister := ps.getProcessToRegister(user, opts, scope)

	existingProcess, err := ps.processRepository.GetByID(ctx, scope, processToRegister.ID)
	if err != nil && !errors.Is(err, ErrRegisteredProcessNotFound) {
		return nil, err
	}

	if existingProcess != nil {
		if !ps.canProcessBeUpdated(existingProcess) {
			return nil, ErrProcessAlreadyRegistered
		}

		err = ps.processRepository.Update(ctx, scope, processToRegister)
		if err != nil {
			return nil, fmt.Errorf("updating registered process: %w", err)
		}
	} else {
		err = ps.processRepository.Create(ctx, scope, processToRegister)
		if err != nil {
			return nil, fmt.Errorf("saving registered process in db: %w", err)
		}
	}

	go ps.uploadProcessToRegistry(scope, processToRegister, opts.Sources)

	return processToRegister, nil
}

func (ps *Service) uploadProcessToRegistry(
	product string,
	registeredProcess *entity.RegisteredProcess,
	sources io.Reader,
) {
	ctx := context.Background()

	ps.logger.Info("Building image")

	tmpFile, err := os.CreateTemp("", "process-compress-*.tar.gz")
	if err != nil {
		ps.uploadingProcessError(ctx, product, registeredProcess, fmt.Errorf("creating temp file for process: %w", err))
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, sources)
	if err != nil {
		ps.uploadingProcessError(ctx, product, registeredProcess, fmt.Errorf("copying temp file for version: %w", err))
		return
	}

	compressedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		ps.uploadingProcessError(ctx, product, registeredProcess, fmt.Errorf("opening process compressed file: %w", err))
		return
	}

	err = ps.objectStorage.UploadImageSources(ctx, product, registeredProcess.Image, compressedFile)
	if err != nil {
		ps.uploadingProcessError(ctx, product, registeredProcess, fmt.Errorf("uploading sources: %w", err))
		return
	}

	defer func() {
		if err := ps.objectStorage.DeleteImageSources(ctx, product, registeredProcess.Image); err != nil {
			ps.logger.Error(err, "Error deleting image's sources", "product", product, "image", registeredProcess.Image)
		}
	}()

	_, err = ps.versionService.RegisterProcess(ctx, product, registeredProcess.ID, registeredProcess.Image)
	if err != nil {
		ps.uploadingProcessError(ctx, product, registeredProcess, fmt.Errorf("registering process: %w", err))
		return
	}

	registeredProcess.Status = entity.RegisterProcessStatusCreated

	err = ps.processRepository.Update(ctx, product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "error updating registered process")

		registeredProcess.Status = entity.RegisterProcessStatusFailed
		registeredProcess.Logs = err.Error()
	}

	ps.logger.Info("Process successfully registered", "processID", registeredProcess.ID)
}

func (ps *Service) uploadingProcessError(
	ctx context.Context,
	product string,
	registeredProcess *entity.RegisteredProcess,
	registerError error,
) {
	ps.logger.Error(registerError, "Error uploading process to registry", "process ID", registeredProcess.ID)
	registeredProcess.Status = entity.RegisterProcessStatusFailed
	registeredProcess.Logs = registerError.Error()

	err := ps.processRepository.Update(ctx, product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "Error updating registered process", "process ID", registeredProcess.ID)
	}
}

func (ps *Service) checkGrants(user *entity.User, opts RegisterProcessOpts) error {
	if opts.IsPublic {
		return ps.accessControl.CheckRoleGrants(user, auth.ActRegisterPublicProcess)
	}

	return ps.accessControl.CheckProductGrants(user, opts.Product, auth.ActRegisterProcess)
}

func (ps *Service) getProcessRegisterScope(opts RegisterProcessOpts) string {
	if opts.IsPublic {
		return viper.GetString(config.GlobalRegistryKey)
	}

	return opts.Product
}

func (ps *Service) canProcessBeUpdated(existingProcess *entity.RegisteredProcess) bool {
	return existingProcess.Version == "latest" || existingProcess.Status == entity.RegisterProcessStatusFailed
}

func (ps *Service) getProcessToRegister(user *entity.User, opts RegisterProcessOpts, scope string) *entity.RegisteredProcess {
	processID := ps.getProcessID(scope, opts.Process, opts.Version)

	return &entity.RegisteredProcess{
		ID:         processID,
		Name:       opts.Process,
		Version:    opts.Version,
		Type:       opts.ProcessType,
		Image:      ps.getProcessImage(processID),
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      user.Email,
		Status:     entity.RegisterProcessStatusCreating,
		IsPublic:   opts.IsPublic,
	}
}

func (ps *Service) getProcessID(scope, process, version string) string {
	return fmt.Sprintf("%s_%s:%s", scope, process, version)
}

func (ps *Service) getProcessImage(processID string) string {
	return fmt.Sprintf("%s/%s", viper.GetString(config.RegistryHostKey), processID)
}
