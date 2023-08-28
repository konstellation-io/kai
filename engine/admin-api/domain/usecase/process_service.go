package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

var (
	ErrInvalidProcessReference = errors.New("invalid process reference")
)

type ProcessService struct {
	logger                    logr.Logger
	processRegistryRepository repository.ProcessRegistryRepo
	processRegistry           service.K8sService
}

type ProcessMetadata struct {
	Dockerfile string
}

func NewProcessService(
	logger logr.Logger,
	k8sService service.K8sService,
	processRegistryRepository repository.ProcessRegistryRepo,
) *ProcessService {
	return &ProcessService{
		logger:                    logger,
		processRegistry:           k8sService,
		processRegistryRepository: processRegistryRepository,
	}
}

func (ps *ProcessService) RegisterProcess(
	ctx context.Context,
	user *entity.User,
	product, version, process, processType string,
	sources io.Reader,
) (string, error) {
	ps.logger.Info("Registering new process")

	tmpFile, err := os.CreateTemp("", "process-compress-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("creating temp file for process: %w", err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, sources)
	if err != nil {
		return "", fmt.Errorf("copying temp file for version: %w", err)
	}

	compressedFile, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("opening process compressed file: %w", err)
	}

	processRef, err := ps.processRegistry.RegisterProcess(ctx, product, version, process, compressedFile)
	if err != nil {
		return "", fmt.Errorf("registering process: %w", err)
	}

	processID, err := ps.getProcessID(processRef)
	if err != nil {
		return "", err
	}

	registeredProcess := &entity.ProcessRegistry{
		ID:         processID,
		Name:       process,
		Version:    version,
		Type:       processType,
		Image:      processRef,
		UploadDate: time.Now(),
		Owner:      user.ID,
	}

	_, err = ps.processRegistryRepository.Create(product, registeredProcess)
	if err != nil {
		return "", fmt.Errorf("saving process registry in db: %w", err)
	}

	ps.logger.Info("Registered process", "processRef", processRef)

	return processID, nil
}

func (ps *ProcessService) getProcessID(processRef string) (string, error) {
	split := strings.Split(processRef, "/")
	if len(split) != 2 {
		return "", ErrInvalidProcessReference
	}

	return split[1], nil
}
