package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
	"github.com/spf13/viper"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

var (
	ErrInvalidProcessType        = errors.New("invalid process type")
	ErrRegisteredProcessNotFound = errors.New("registered process not found")
	ErrProcessAlreadyRegistered  = errors.New("process already registered")
)

type ProcessService struct {
	logger            logr.Logger
	processRepository repository.ProcessRepository
	versionService    service.VersionService
}

type ProcessMetadata struct {
	Dockerfile string
}

func NewProcessService(
	logger logr.Logger,
	k8sService service.VersionService,
	processRepository repository.ProcessRepository,
) *ProcessService {
	return &ProcessService{
		logger:            logger,
		versionService:    k8sService,
		processRepository: processRepository,
	}
}

func (ps *ProcessService) RegisterProcess(
	ctx context.Context,
	user *entity.User,
	product, version, process, processType string,
	sources io.Reader,
) (*entity.RegisteredProcess, chan *entity.RegisteredProcess, error) {
	ps.logger.Info("Registering new process")

	processID := fmt.Sprintf("%s_%s:%s", product, process, version)

	existingProcess, err := ps.processRepository.GetByID(ctx, product, processID)
	if err != nil && !errors.Is(err, ErrRegisteredProcessNotFound) {
		return nil, nil, err
	}

	registryURL, err := url.Parse(viper.GetString(config.RegistryURLKey))
	if err != nil {
		return nil, nil, fmt.Errorf("parsing registry url: %w", err)
	}

	processImage := fmt.Sprintf("%s/%s", registryURL.Host, processID)

	registeredProcess := &entity.RegisteredProcess{
		ID:         processID,
		Name:       process,
		Version:    version,
		Type:       processType,
		Image:      processImage,
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      user.ID,
		Status:     entity.RegisterProcessStatusCreating,
	}

	processExists := existingProcess != nil
	if processExists {
		processStatusIsFailed := existingProcess.Status == entity.RegisterProcessStatusFailed
		if !processStatusIsFailed {
			return nil, nil, ErrProcessAlreadyRegistered
		}

		err = ps.processRepository.Update(ctx, product, registeredProcess)
		if err != nil {
			return nil, nil, fmt.Errorf("updating registered process: %w", err)
		}
	} else {
		_, err = ps.processRepository.Create(product, registeredProcess)
		if err != nil {
			return nil, nil, fmt.Errorf("saving process registry in db: %w", err)
		}
	}

	notifyStatusCh := make(chan *entity.RegisteredProcess, 1)

	go ps.uploadProcessToRegistry(product, registeredProcess, sources, notifyStatusCh)

	return registeredProcess, notifyStatusCh, nil
}

func (ps *ProcessService) uploadProcessToRegistry(
	product string,
	registeredProcess *entity.RegisteredProcess,
	sources io.Reader,
	notifyStatusCh chan *entity.RegisteredProcess,
) {
	ctx := context.Background()

	defer close(notifyStatusCh)

	ps.logger.Info("Building image")

	tmpFile, err := os.CreateTemp("", "process-compress-*.tar.gz")
	if err != nil {
		ps.updatedRegisteredProcessError(ctx, product, registeredProcess, notifyStatusCh, fmt.Errorf("creating temp file for process: %w", err))
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, sources)
	if err != nil {
		ps.updatedRegisteredProcessError(ctx, product, registeredProcess, notifyStatusCh, fmt.Errorf("copying temp file for version: %w", err))
		return
	}

	compressedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		ps.updatedRegisteredProcessError(ctx, product, registeredProcess, notifyStatusCh, fmt.Errorf("opening process compressed file: %w", err))
		return
	}

	_, err = ps.versionService.RegisterProcess(ctx, registeredProcess.ID, registeredProcess.Image, compressedFile)
	if err != nil {
		ps.updatedRegisteredProcessError(ctx, product, registeredProcess, notifyStatusCh, fmt.Errorf("registering process: %w", err))
		return
	}

	registeredProcess.Status = entity.RegisterProcessStatusCreated

	err = ps.processRepository.Update(ctx, product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "Error updating registered process")
		registeredProcess.Status = entity.RegisterProcessStatusFailed
		registeredProcess.Logs = err.Error()
	}

	notifyStatusCh <- registeredProcess
}

func (ps *ProcessService) updatedRegisteredProcessError(
	ctx context.Context,
	product string,
	registeredProcess *entity.RegisteredProcess,
	notifyStatusCh chan *entity.RegisteredProcess,
	registerError error,
) {
	registeredProcess.Status = entity.RegisterProcessStatusFailed
	registeredProcess.Logs = registerError.Error()

	err := ps.processRepository.Update(ctx, product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "Error updating registered process", "process ID", registeredProcess.ID)
	}

	notifyStatusCh <- registeredProcess
}

// TODO: move this to viper config
//
//nolint:godox // To be done.
func (ps *ProcessService) ListByProductAndType(
	ctx context.Context,
	user *entity.User,
	productID, processType string,
) ([]*entity.RegisteredProcess, error) {
	log := fmt.Sprintf("Retrieving process for product %q", productID)
	if processType != "" {
		log = fmt.Sprintf("%s with process type filter %q", log, processType)
	}

	ps.logger.Info(log)

	if processType != "" && !entity.ProcessType(processType).IsValid() {
		return nil, fmt.Errorf("%w: %s", ErrInvalidProcessType, processType)
	}

	return ps.processRepository.ListByProductAndType(ctx, productID, processType)
}
