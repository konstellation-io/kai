package usecase

import (
	"context"
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

type ProcessService struct {
	logger logr.Logger
	//processRegistry   ProcessRegistry
	processRegistryRepository repository.ProcessRegistryRepo
	processRegistry           service.K8sService
}

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/${GOFILE} -package=mocks

//type ProcessRegistry interface {
//	RegisterProcess(ctx context.Context, product, version, process string, src io.Reader) (string, error)
//}

type ProcessMetadata struct {
	Dockerfile string
}

func NewProcessService(
	logger logr.Logger,
	//processRegistry ProcessRegistry,
	k8sService service.K8sService,
	processRegistryRepository repository.ProcessRegistryRepo,
) *ProcessService {
	return &ProcessService{
		logger: logger,
		//processRegistry: processRegistry,
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
	//defer os.Remove(tmpFile.Name())

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

	processID := ps.getProcessID(processRef)
	registeredProcess := &entity.ProcessRegistry{
		ID:         processID,
		Name:       process,
		Version:    version,
		Type:       processType,
		Image:      processRef,
		UploadDate: time.Now().Truncate(time.Millisecond).UTC(),
		Owner:      user.ID,
	}

	_, err = ps.processRegistryRepository.Create(product, registeredProcess)
	if err != nil {
		return "", fmt.Errorf("saving process registry in db: %w", err)
	}

	ps.logger.Info("Registered process", "processRef", processRef)

	return processID, nil
}

func (ps *ProcessService) getProcessID(processRef string) string {
	split := strings.Split(processRef, "/")
	if len(split) != 2 {
		ps.logger.Error(fmt.Errorf("invalid process ref: %s", processRef), "defaulting to use whole processRef")
		return processRef
	}
	return split[1]
}

func (ps *ProcessService) ListByProductWithTypeFilter(
	ctx context.Context,
	user *entity.User,
	productID, processType string,
) ([]*entity.ProcessRegistry, error) {
	ps.logger.Info("Retrieving process for product %q with process type filter %q", productID, processType)

	if !entity.ProcessType(processType).IsValid() {
		return nil, fmt.Errorf("invalid process type: %s", processType)
	}

	return ps.processRegistryRepository.ListByProductWithTypeFilter(ctx, productID, processType)
}
