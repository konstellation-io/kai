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
) (*entity.RegisteredProcess, error) {
	ps.logger.Info("Registering new process")

	processID := fmt.Sprintf("%s_%s:%s", product, process, version)

	existingProcess, err := ps.processRepository.GetByID(ctx, product, processID)
	if err != nil && !errors.Is(err, ErrRegisteredProcessNotFound) {
		return nil, err
	}

	registryURL, err := url.Parse(viper.GetString(config.RegistryURLKey))
	if err != nil {
		return nil, fmt.Errorf("parsing registry url: %w", err)
	}

	//return fmt.Sprintf("%s/%s", registryURL.Host, imageName), nil
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
		processStatusIsFailed := existingProcess.Status != entity.RegisterProcessStatusFailed
		if processStatusIsFailed {
			return nil, ErrProcessAlreadyRegistered
		}

		err = ps.processRepository.Update(product, registeredProcess)
		if err != nil {
			return nil, fmt.Errorf("updating registered process: %w", err)
		}

	} else {
		_, err = ps.processRepository.Create(product, registeredProcess)
		if err != nil {
			return nil, fmt.Errorf("saving process registry in db: %w", err)
		}
	}

	go ps.uploadProcessToRegistry(product, registeredProcess, sources)

	return registeredProcess, nil
}

func (ps *ProcessService) uploadProcessToRegistry(product string, registeredProcess *entity.RegisteredProcess, sources io.Reader) {
	ctx := context.Background()
	ps.logger.Info("Building image")

	tmpFile, err := os.CreateTemp("", "process-compress-*.tar.gz")
	if err != nil {
		ps.updatedRegisteredProcessError(product, registeredProcess, fmt.Errorf("creating temp file for process: %s", err))
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, sources)
	if err != nil {
		ps.updatedRegisteredProcessError(product, registeredProcess, fmt.Errorf("copying temp file for version: %w", err))
		return
	}

	compressedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		ps.updatedRegisteredProcessError(product, registeredProcess, fmt.Errorf("opening process compressed file: %w", err))
		return
	}

	_, err = ps.versionService.RegisterProcess(ctx, registeredProcess.ID, registeredProcess.Image, compressedFile)
	if err != nil {
		ps.updatedRegisteredProcessError(product, registeredProcess, fmt.Errorf("registering process: %w", err))
		return
	}

	registeredProcess.Status = entity.RegisterProcessStatusCreated

	err = ps.processRepository.Update(product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "Error updating registered process")
		return
	}

	ps.logger.Info("Process successfully registered", "processID", registeredProcess.ID)
}

func (ps *ProcessService) updatedRegisteredProcessError(
	product string,
	registeredProcess *entity.RegisteredProcess,
	registerError error,
) {
	registeredProcess.Status = entity.RegisterProcessStatusFailed
	registeredProcess.Logs = registerError.Error()

	err := ps.processRepository.Update(product, registeredProcess)
	if err != nil {
		ps.logger.Error(err, "Error updating registered process", "process ID", registeredProcess.ID)
		return
	}
	return
}

func (ps *ProcessService) ListByProductAndType(
	ctx context.Context,
	user *entity.User,
	productID, processType string,
) ([]*entity.RegisteredProcess, error) {
	// TODO: Check user grants

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
