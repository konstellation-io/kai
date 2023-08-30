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
	ErrInvalidProcessType      = errors.New("invalid process type")
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

	processRef, err := ps.versionService.RegisterProcess(ctx, product, version, process, compressedFile)
	if err != nil {
		return "", fmt.Errorf("registering process: %w", err)
	}

	processID := ps.getProcessID(processRef)
	registeredProcess := &entity.RegisteredProcess{
		ID:         processID,
		Name:       process,
		Version:    version,
		Type:       processType,
		Image:      processRef,
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      user.ID,
	}

	_, err = ps.processRepository.Create(product, registeredProcess)
	if err != nil {
		return "", fmt.Errorf("saving process registry in db: %w", err)
	}

	ps.logger.Info("Registered process with ID:", processID)

	return processID, nil
}

func (ps *ProcessService) getProcessID(processRef string) string {
	split := strings.Split(processRef, "/")

	if len(split) != 2 {
		ps.logger.Info(
			fmt.Sprintf("WARNING: invalid process reference %q, defaulting to use whole process reference", processRef),
		)
		return processRef
	}

	return split[1]
}

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
