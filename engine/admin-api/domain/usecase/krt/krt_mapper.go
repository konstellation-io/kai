package krt

import (
	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// KrtMapper is a service to map KRT YAML to entity.Version
//
// caution: this method takes for granted the KRT YAML is valid
func MapKrtYamlToVersion(krtYml *krt.Krt) *entity.Version {
	version := &entity.Version{
		Name:        krtYml.Name,
		Description: krtYml.Description,
		Version:     krtYml.Version,
		Config:      keyValueMapToConfigurationVariableArray(krtYml.Config),
		Workflows:   mapKrtYamlToWorkflows(krtYml.Workflows),
	}

	return version
}

func mapKrtYamlToWorkflows(krtWorkflows []krt.Workflow) []entity.Workflow {
	workflows := make([]entity.Workflow, len(krtWorkflows))

	for i, krtWorkflow := range krtWorkflows {
		workflows[i] = entity.Workflow{
			Name:      krtWorkflow.Name,
			Type:      entity.WorkflowType(krtWorkflow.Type),
			Config:    keyValueMapToConfigurationVariableArray(krtWorkflow.Config),
			Processes: mapKrtYamlToProcesses(krtWorkflow.Processes),
		}
	}

	return workflows
}

func mapKrtYamlToProcesses(krtProcesses []krt.Process) []entity.Process {
	processes := make([]entity.Process, len(krtProcesses))

	for i, krtProcess := range krtProcesses {
		processes[i] = entity.Process{
			Name:          krtProcess.Name,
			Type:          entity.ProcessType(krtProcess.Type),
			Image:         krtProcess.Image,
			Replicas:      int32(*krtProcess.Replicas),
			GPU:           *krtProcess.GPU,
			Config:        keyValueMapToConfigurationVariableArray(krtProcess.Config),
			ObjectStore:   mapKrtYamlToProcessObjectStore(krtProcess.ObjectStore),
			Secrets:       keyValueMapToConfigurationVariableArray(krtProcess.Secrets),
			Subscriptions: krtProcess.Subscriptions,
			Networking:    mapKrtYamlToProcessNetworking(krtProcess.Networking),
		}
	}

	return processes
}

func mapKrtYamlToProcessObjectStore(krtObjectStore *krt.ProcessObjectStore) *entity.ProcessObjectStore {
	if krtObjectStore == nil {
		return nil
	}

	return &entity.ProcessObjectStore{
		Name:  krtObjectStore.Name,
		Scope: entity.ObjectStoreScope(krtObjectStore.Scope),
	}
}

func mapKrtYamlToProcessNetworking(krtNetworking *krt.ProcessNetworking) *entity.ProcessNetworking {
	if krtNetworking == nil {
		return nil
	}

	return &entity.ProcessNetworking{
		TargetPort:      krtNetworking.TargetPort,
		DestinationPort: krtNetworking.DestinationPort,
		Protocol:        string(krtNetworking.Protocol),
	}
}

func keyValueMapToConfigurationVariableArray(m map[string]string) []entity.ConfigurationVariable {
	if m == nil {
		return nil
	}

	config := make([]entity.ConfigurationVariable, len(m))
	idx := 0

	for k, v := range m {
		config[idx] = entity.ConfigurationVariable{
			Key:   k,
			Value: v,
		}
		idx++
	}

	return config
}
