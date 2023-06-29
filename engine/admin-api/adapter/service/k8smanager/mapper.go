package k8smanager

import (
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapWorkflowsToDTO(workflows []entity.Workflow, versionConfig *entity.VersionConfig) ([]*versionpb.Workflow, error) {
	workflowsDTO := make([]*versionpb.Workflow, 0, len(workflows))

	for _, w := range workflows {
		wStreamCfg, err := versionConfig.GetWorkflowStreamConfig(w.Name)
		if err != nil {
			return nil, fmt.Errorf("get workflows's %q stream config: %w", w.Name, err)
		}

		wKeyValueCfg, err := versionConfig.GetWorkflowKeyValueStoresConfig(w.Name)
		if err != nil {
			return nil, fmt.Errorf("get workflow's %q key-value store config: %w", w.Name, err)
		}

		wObjectStoreCfg, err := versionConfig.GetWorkflowObjectStoresConfig(w.Name)
		if err != nil {
			return nil, fmt.Errorf("get worklfow's %q object store config: %w", w.Name, err)
		}

		processesDTO, err := mapProcessesToDTO(w.Processes, wStreamCfg, wKeyValueCfg, wObjectStoreCfg)
		if err != nil {
			return nil, fmt.Errorf("map workflow's %q process to dto: %w", w.Name, err)
		}

		workflowsDTO = append(workflowsDTO, &versionpb.Workflow{
			Id:            w.ID,
			Processes:     processesDTO,
			Stream:        wStreamCfg.Stream,
			KeyValueStore: wKeyValueCfg.WorkflowKeyValueStore,
		})
	}

	return workflowsDTO, nil
}

func mapProcessesToDTO(
	processes []entity.Process,
	streamConfig *entity.WorkflowStreamConfig,
	kvConfig *entity.WorkflowKeyValueStores,
	objStoreConfig *entity.WorkflowObjectStoresConfig,
) ([]*versionpb.Process, error) {
	processesDTO := make([]*versionpb.Process, 0, len(processes))

	for _, p := range processes {
		nodeStreamCfg, err := streamConfig.GetProcessConfig(p.Name)
		if err != nil {
			return nil, fmt.Errorf("get node's %q stream config: %w", p.Name, err)
		}

		keyValueStore, err := kvConfig.GetProcessKeyValueStore(p.Name)
		if err != nil {
			return nil, fmt.Errorf("get process' %q key value store config: %w", p.Name, err)
		}

		process := &versionpb.Process{
			Id:            p.ID,
			Image:         p.Image,
			Gpu:           p.GPU,
			Subscriptions: nodeStreamCfg.Subscriptions,
			Subject:       nodeStreamCfg.Subject,
			KeyValueStore: keyValueStore,
			Replicas:      p.Replicas,
			Config:        mapProcessConfigToDTO(p.Config),
			Type:          mapProcessTypeToDTO(p.Type),
		}

		processObjectStore := objStoreConfig.Processes.GetProcessObjectStoreConfig(p.Name)
		if processObjectStore != nil {
			process.ObjectStore = processObjectStore
		}

		if p.Networking != nil {
			process.Networking = &versionpb.Network{
				TargetPort: int32(p.Networking.TargetPort),
				Protocol:   p.Networking.Protocol,
				SourcePort: int32(p.Networking.DestinationPort),
			}
		}

		processesDTO = append(processesDTO, process)
	}

	return processesDTO, nil
}

func mapProcessConfigToDTO(config []entity.ConfigurationVariable) map[string]string {
	if len(config) == 0 {
		return nil
	}

	configuration := make(map[string]string, len(config))

	for _, c := range config {
		configuration[c.Key] = c.Value
	}

	return configuration
}

func mapProcessTypeToDTO(processType entity.ProcessType) versionpb.ProcessType {
	switch processType {
	case entity.ProcessTypeTrigger:
		return versionpb.ProcessType_ProcessTypeTrigger
	case entity.ProcessTypeTask:
		return versionpb.ProcessType_ProcessTypeTask
	case entity.ProcessTypeExit:
		return versionpb.ProcessType_ProcessTypeExit
	default:
		return versionpb.ProcessType_ProcessTypeUnknown
	}
}
