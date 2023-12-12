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

	var registry string

	if opts.IsPublic {
		if err := ps.accessControl.CheckRoleGrants(user, auth.ActRegisterPublicProcess); err != nil {
			return nil, err
		}

		registry = viper.GetString(config.GlobalRegistryKey)
	} else {
		if err := ps.accessControl.CheckProductGrants(user, opts.Product, auth.ActRegisterProcess); err != nil {
			return nil, err
		}

		registry = opts.Product
	}

	processID := fmt.Sprintf("%s_%s:%s", registry, opts.Process, opts.Version)

	existingProcess, err := ps.processRepository.GetByID(ctx, registry, processID)
	if err != nil && !errors.Is(err, ErrRegisteredProcessNotFound) {
		return nil, err
	}

	processImage := fmt.Sprintf("%s/%s", viper.GetString(config.RegistryHostKey), processID)

	registeredProcess := &entity.RegisteredProcess{
		ID:         processID,
		Name:       opts.Process,
		Version:    opts.Version,
		Type:       opts.ProcessType,
		Image:      processImage,
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      user.Email,
		Status:     entity.RegisterProcessStatusCreating,
	}

	processExists := existingProcess != nil

	if processExists {
		isLatest := opts.Version == "latest"
		processStatusIsFailed := existingProcess.Status == entity.RegisterProcessStatusFailed

		if !processStatusIsFailed && !isLatest {
			return nil, ErrProcessAlreadyRegistered
		}

		err = ps.processRepository.Update(ctx, registry, registeredProcess)
		if err != nil {
			return nil, fmt.Errorf("updating registered process: %w", err)
		}
	} else {
		err = ps.processRepository.Create(ctx, registry, registeredProcess)
		if err != nil {
			return nil, fmt.Errorf("saving process registry in db: %w", err)
		}
	}

	go ps.uploadProcessToRegistry(registry, registeredProcess, opts.Sources)

	return registeredProcess, nil
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
