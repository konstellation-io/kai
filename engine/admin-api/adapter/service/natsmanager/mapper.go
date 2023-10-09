package natsmanager

import (
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func (n *Client) mapWorkflowsToDTO(workflows []entity.Workflow) []*natspb.Workflow {
	workflowsDTO := make([]*natspb.Workflow, 0, len(workflows))

	for _, w := range workflows {
		processes := make([]*natspb.Process, 0, len(w.Processes))

		for _, process := range w.Processes {
			processToAppend := natspb.Process{
				Name:          process.Name,
				Subscriptions: process.Subscriptions,
			}

			if process.ObjectStore != nil {
				processToAppend.ObjectStore = &natspb.ObjectStore{
					Name:  process.ObjectStore.Name,
					Scope: mapObjectStoreScopeToDTO(process.ObjectStore.Scope),
				}
			}

			processes = append(processes, &processToAppend)
		}

		workflowsDTO = append(workflowsDTO, &natspb.Workflow{
			Name:      w.Name,
			Processes: processes,
		})
	}

	return workflowsDTO
}

func (n *Client) mapDTOToVersionStreamConfig(
	workflowsDTO map[string]*natspb.WorkflowStreamConfig,
) *entity.VersionStreamsConfig {
	workflows := make(map[string]entity.WorkflowStreamConfig, len(workflowsDTO))

	for workflow, streamCfg := range workflowsDTO {
		workflows[workflow] = entity.WorkflowStreamConfig{
			Stream:    streamCfg.Stream,
			Processes: n.mapDTOToProcessesStreamConfig(streamCfg.Processes),
		}
	}

	return &entity.VersionStreamsConfig{
		Workflows: workflows,
	}
}

func (n *Client) mapDTOToProcessesStreamConfig(
	processes map[string]*natspb.ProcessStreamConfig,
) map[string]entity.ProcessStreamConfig {
	processesStreamCfg := map[string]entity.ProcessStreamConfig{}

	for process, subjectCfg := range processes {
		processesStreamCfg[process] = entity.ProcessStreamConfig{
			Subject:       subjectCfg.Subject,
			Subscriptions: subjectCfg.Subscriptions,
		}
	}

	return processesStreamCfg
}

func (n *Client) mapDTOToVersionObjectStoreConfig(
	workflowsDTO map[string]*natspb.WorkflowObjectStoreConfig,
) *entity.VersionObjectStoresConfig {
	workflows := make(map[string]entity.WorkflowObjectStoresConfig, len(workflowsDTO))

	for workflow, objStoreCfg := range workflowsDTO {
		workflows[workflow] = entity.WorkflowObjectStoresConfig{
			Processes: objStoreCfg.Processes,
		}
	}

	return &entity.VersionObjectStoresConfig{
		Workflows: workflows,
	}
}

func (n *Client) mapDTOToVersionKeyValueStoreConfig(
	projectKeyValueStore string,
	workflows map[string]*natspb.WorkflowKeyValueStoreConfig,
) *entity.KeyValueStores {
	workflowsKVConfig := make(map[string]*entity.WorkflowKeyValueStores, len(workflows))

	for workflow, kvStoreCfg := range workflows {
		workflowsKVConfig[workflow] = &entity.WorkflowKeyValueStores{
			KeyValueStore: kvStoreCfg.KeyValueStore,
			Processes:     kvStoreCfg.Processes,
		}
	}

	return &entity.KeyValueStores{
		KeyValueStore: projectKeyValueStore,
		Workflows:     workflowsKVConfig,
	}
}

func mapObjectStoreScopeToDTO(scope entity.ObjectStoreScope) natspb.ObjectStoreScope {
	//nolint:exhaustive // wrong lint rule
	switch scope {
	case "project":
		return natspb.ObjectStoreScope_SCOPE_PROJECT
	case "workflow":
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW
	default:
		return natspb.ObjectStoreScope_SCOPE_UNDEFINED
	}
}
