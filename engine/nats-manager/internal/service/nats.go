package service

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/logging"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

// NatsService basic server.
type NatsService struct {
	config  *config.Config
	logger  logging.Logger
	manager *manager.NatsManager
	natspb.UnimplementedNatsManagerServiceServer
}

// NewNatsService instantiates the GRPC server implementation.
func NewNatsService(
	cfg *config.Config,
	logger logging.Logger,
	manager *manager.NatsManager,
) *NatsService {
	return &NatsService{
		cfg,
		logger,
		manager,
		natspb.UnimplementedNatsManagerServiceServer{},
	}
}

// CreateStreams create streams for given workflows.
func (n *NatsService) CreateStreams(
	_ context.Context,
	req *natspb.CreateStreamsRequest,
) (*natspb.CreateStreamsResponse, error) {
	n.logger.Info("CreateStreams request received")

	workflows := n.dtoToWorkflows(req.Workflows)

	streamConfig, err := n.manager.CreateStreams(req.ProductId, req.VersionName, workflows)
	if err != nil {
		n.logger.Errorf("Error creating streams: %s", err)
		return nil, err
	}

	return &natspb.CreateStreamsResponse{
		Workflows: n.workflowsStreamConfigToDto(streamConfig),
	}, nil
}

// CreateObjectStores creates object stores for given workflows.
func (n *NatsService) CreateObjectStores(
	_ context.Context,
	req *natspb.CreateObjectStoresRequest,
) (*natspb.CreateObjectStoresResponse, error) {
	n.logger.Info("CreateObjectStores request received")

	objectStores, err := n.manager.CreateObjectStores(req.ProductId, req.VersionName, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Errorf("Error creating object store: %s", err)
		return nil, err
	}
	return &natspb.CreateObjectStoresResponse{
		Workflows: n.mapWorkflowsObjStoreToDTO(objectStores),
	}, nil
}

// DeleteStreams delete streams for given workflows.
func (n *NatsService) DeleteStreams(
	_ context.Context,
	req *natspb.DeleteStreamsRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("Delete streams request received")

	err := n.manager.DeleteStreams(req.ProductId, req.VersionName)
	if err != nil {
		n.logger.Errorf("Error deleting streams: %s", err)
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Streams and subjects for version %q on product %s deleted", req.VersionName, req.ProductId),
	}, nil
}

// DeleteObjectStores delete object stores for given workflows.
func (n *NatsService) DeleteObjectStores(
	_ context.Context,
	req *natspb.DeleteObjectStoresRequest,
) (*natspb.DeleteResponse, error) {
	n.logger.Info("Delete object stores request received")

	err := n.manager.DeleteObjectStores(req.ProductId, req.VersionName)
	if err != nil {
		n.logger.Errorf("Error deleting object stores: %s", err)
		return nil, err
	}

	return &natspb.DeleteResponse{
		Message: fmt.Sprintf("Object stores for version %q on product %s deleted", req.VersionName, req.ProductId),
	}, nil
}

func (n *NatsService) CreateKeyValueStores(
	_ context.Context,
	req *natspb.CreateKeyValueStoresRequest,
) (*natspb.CreateKeyValueStoreResponse, error) {
	n.logger.Info("CreateKeyValueStores request received")

	keyValueStores, err := n.manager.CreateKeyValueStores(req.ProductId, req.VersionName, n.dtoToWorkflows(req.Workflows))
	if err != nil {
		n.logger.Errorf("Error creating object store: %s", err)
		return nil, err
	}
	return n.mapKeyValueStoresToDTO(keyValueStores), nil
}

func (n *NatsService) dtoToWorkflows(dtoWorkflows []*natspb.Workflow) []entity.Workflow {
	workflows := make([]entity.Workflow, 0, len(dtoWorkflows))

	for _, dtoWorkflow := range dtoWorkflows {
		workflows = append(workflows, entity.Workflow{
			Name:      dtoWorkflow.Name,
			Processes: n.dtoToProcesses(dtoWorkflow.Processes),
		})
	}

	return workflows
}

func (n *NatsService) dtoToProcesses(processesDTO []*natspb.Process) []entity.Process {
	processes := make([]entity.Process, 0, len(processesDTO))

	for _, processDTO := range processesDTO {
		process := entity.Process{
			Name:          processDTO.Name,
			Subscriptions: processDTO.Subscriptions,
		}

		if processDTO.ObjectStore != nil {
			process.ObjectStore = &entity.ObjectStore{
				Name:  processDTO.ObjectStore.Name,
				Scope: entity.ObjectStoreScope(processDTO.ObjectStore.Scope),
			}
		}
		processes = append(processes, process)
	}

	return processes
}

func (n *NatsService) workflowsStreamConfigToDto(
	workflows entity.WorkflowsStreamsConfig,
) map[string]*natspb.WorkflowStreamConfig {
	workflowsStreamCfg := make(map[string]*natspb.WorkflowStreamConfig, len(workflows))

	for workflow, cfg := range workflows {
		workflowsStreamCfg[workflow] = &natspb.WorkflowStreamConfig{
			Stream:    cfg.Stream,
			Processes: n.mapProcessStreamConfigToDTO(cfg.Processes),
		}
	}

	return workflowsStreamCfg
}

func (n *NatsService) mapProcessStreamConfigToDTO(
	processes entity.ProcessesStreamConfig,
) map[string]*natspb.ProcessStreamConfig {
	processesStreamCfg := make(map[string]*natspb.ProcessStreamConfig, len(processes))

	for process, cfg := range processes {
		processesStreamCfg[process] = &natspb.ProcessStreamConfig{
			Subject:       cfg.Subject,
			Subscriptions: cfg.Subscriptions,
		}
	}

	return processesStreamCfg
}

func (n *NatsService) mapWorkflowsObjStoreToDTO(
	workflowsObjStores entity.WorkflowsObjectStoresConfig,
) map[string]*natspb.WorkflowObjectStoreConfig {
	workflowsConfig := map[string]*natspb.WorkflowObjectStoreConfig{}

	for workflow, objectStoresConfig := range workflowsObjStores {
		workflowsConfig[workflow] = &natspb.WorkflowObjectStoreConfig{
			Processes: objectStoresConfig.Processes,
		}
	}

	return workflowsConfig
}

func (n *NatsService) mapKeyValueStoresToDTO(stores *entity.VersionKeyValueStores) *natspb.CreateKeyValueStoreResponse {
	workflowsStores := make(map[string]*natspb.WorkflowKeyValueStoreConfig, len(stores.WorkflowsStores))

	for workflow, storesConfig := range stores.WorkflowsStores {
		workflowsStores[workflow] = &natspb.WorkflowKeyValueStoreConfig{
			KeyValueStore: storesConfig.WorkflowStore,
			Processes:     storesConfig.Processes,
		}
	}

	return &natspb.CreateKeyValueStoreResponse{
		KeyValueStore: stores.ProjectStore,
		Workflows:     workflowsStores,
	}
}
