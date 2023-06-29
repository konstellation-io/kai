package natsmanager

import (
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func (n *NatsManagerClient) mapWorkflowsToDTO(workflows []entity.Workflow) ([]*natspb.Workflow, error) {
	workflowsDTO := make([]*natspb.Workflow, 0, len(workflows))

	for _, w := range workflows {
		processes := make([]*natspb.Process, 0, len(w.Processes))

		for _, process := range w.Processes {
			processToAppend := natspb.Process{
				Id:            process.Name,
				Subscriptions: process.Subscriptions,
			}

			if process.ObjectStore != nil {
				scope, err := mapObjectStoreScopeToDTO(process.ObjectStore.Scope)
				if err != nil {
					return nil, err
				}

				processToAppend.ObjectStore = &natspb.ObjectStore{
					Name:  process.ObjectStore.Name,
					Scope: scope,
				}
			}

			processes = append(processes, &processToAppend)
		}

		workflowsDTO = append(workflowsDTO, &natspb.Workflow{
			Id:        w.Name,
			Processes: processes,
		})
	}

	return workflowsDTO, nil
}

func (n *NatsManagerClient) mapDTOToVersionStreamConfig(
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

func (n *NatsManagerClient) mapDTOToProcessesStreamConfig(
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

func (n *NatsManagerClient) mapDTOToVersionObjectStoreConfig(
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

func (n *NatsManagerClient) mapDTOToVersionKeyValueStoreConfig(
	projectKeyValueStore string,
	workflows map[string]*natspb.WorkflowKeyValueStoreConfig,
) *entity.KeyValueStoresConfig {
	workflowsKVConfig := make(map[string]*entity.WorkflowKeyValueStores, len(workflows))

	for workflow, kvStoreCfg := range workflows {
		workflowsKVConfig[workflow] = &entity.WorkflowKeyValueStores{
			WorkflowKeyValueStore:   kvStoreCfg.KeyValueStore,
			ProcessesKeyValueStores: kvStoreCfg.Processes,
		}
	}

	return &entity.KeyValueStoresConfig{
		ProductKeyValueStore:    projectKeyValueStore,
		WorkflowsKeyValueStores: workflowsKVConfig,
	}
}

func mapObjectStoreScopeToDTO(scope entity.ObjectStoreScope) (natspb.ObjectStoreScope, error) {
	//nolint:exhaustive // wrong lint rule
	switch scope {
	case "project":
		return natspb.ObjectStoreScope_SCOPE_PROJECT, nil
	case "workflow":
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, nil
	default:
		//nolint:goerr113 // error needs to be wrapped
		return natspb.ObjectStoreScope_SCOPE_WORKFLOW, errors.New("invalid object store scope")
	}
}
