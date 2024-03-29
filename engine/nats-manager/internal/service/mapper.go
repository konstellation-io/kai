package service

import (
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
)

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
				Scope: n.dtoToObjectStoreScope(processDTO.ObjectStore.Scope),
			}
		}

		processes = append(processes, process)
	}

	return processes
}

func (n *NatsService) dtoToObjectStoreScope(scope natspb.ObjectStoreScope) entity.ObjectStoreScope {
	switch scope {
	case natspb.ObjectStoreScope_SCOPE_PROJECT:
		return entity.ObjStoreScopeProject
	case natspb.ObjectStoreScope_SCOPE_WORKFLOW:
		return entity.ObjStoreScopeWorkflow
	case natspb.ObjectStoreScope_SCOPE_UNDEFINED:
		return entity.ObjStoreScopeUndefined
	default:
		return entity.ObjStoreScopeUndefined
	}
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

func (n *NatsService) mapKeyValueStoresToDTO(stores *entity.VersionKeyValueStores) *natspb.CreateVersionKeyValueStoresResponse {
	workflowsStores := make(map[string]*natspb.WorkflowKeyValueStoreConfig, len(stores.WorkflowsStores))

	for workflow, storesConfig := range stores.WorkflowsStores {
		workflowsStores[workflow] = &natspb.WorkflowKeyValueStoreConfig{
			KeyValueStore: storesConfig.WorkflowStore,
			Processes:     storesConfig.Processes,
		}
	}

	return &natspb.CreateVersionKeyValueStoresResponse{
		KeyValueStore: stores.ProjectStore,
		Workflows:     workflowsStores,
	}
}

func (n *NatsService) mapDTOToKeyValueStoreConfigurations(dto []*natspb.KeyValueConfiguration) []entity.KeyValueConfiguration {
	keyValueConfigurations := make([]entity.KeyValueConfiguration, 0, len(dto))

	for _, kvCfgDto := range dto {
		keyValueConfigurations = append(keyValueConfigurations, entity.KeyValueConfiguration{
			KeyValueStore: kvCfgDto.KeyValueStore,
			Configuration: kvCfgDto.Configuration,
		})
	}

	return keyValueConfigurations
}
