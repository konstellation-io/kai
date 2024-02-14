package version

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/krt/pkg/krt"
)

func (h *Handler) mapKrtToVersion(krtYaml *krt.Krt) *entity.Version {
	return &entity.Version{
		Tag:         krtYaml.Version,
		Description: krtYaml.Description,
		Config:      h.mapKrtConfigToVersion(krtYaml.Config),
		Workflows:   h.mapKrtWorkflowsToVersion(krtYaml.Workflows),
		Status:      entity.VersionStatusCreated,
	}
}

func (h *Handler) mapKrtConfigToVersion(m map[string]string) []entity.ConfigurationVariable {
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

func (h *Handler) mapKrtWorkflowsToVersion(krtWorkflows []krt.Workflow) []entity.Workflow {
	workflows := make([]entity.Workflow, len(krtWorkflows))

	for i, krtWorkflow := range krtWorkflows {
		workflows[i] = entity.Workflow{
			Name:      krtWorkflow.Name,
			Type:      entity.WorkflowType(krtWorkflow.Type),
			Config:    h.mapKrtConfigToVersion(krtWorkflow.Config),
			Processes: h.mapKrtProcessesToVersion(krtWorkflow.Processes),
		}
	}

	return workflows
}

func (h *Handler) mapKrtProcessesToVersion(krtProcesses []krt.Process) []entity.Process {
	processes := make([]entity.Process, len(krtProcesses))

	for i, krtProcess := range krtProcesses {
		processes[i] = entity.Process{
			Name:           krtProcess.Name,
			Type:           entity.ProcessType(krtProcess.Type),
			Image:          krtProcess.Image,
			Replicas:       int32(*krtProcess.Replicas),
			GPU:            *krtProcess.GPU,
			Config:         h.mapKrtConfigToVersion(krtProcess.Config),
			ObjectStore:    h.mapKrtObjectStoreToVersion(krtProcess.ObjectStore),
			Secrets:        h.mapKrtConfigToVersion(krtProcess.Secrets),
			Subscriptions:  krtProcess.Subscriptions,
			Networking:     h.mapKrtNetworkingToVersion(krtProcess.Networking),
			ResourceLimits: h.mapKrtResourceLimitsToVersion(krtProcess.ResourceLimits),
			Status:         entity.RegisterProcessStatusCreated,
			NodeSelectors:  krtProcess.NodeSelectors,
		}
	}

	return processes
}

func (h *Handler) mapKrtObjectStoreToVersion(krtObjectStore *krt.ProcessObjectStore) *entity.ProcessObjectStore {
	if krtObjectStore == nil {
		return nil
	}

	return &entity.ProcessObjectStore{
		Name:  krtObjectStore.Name,
		Scope: entity.ObjectStoreScope(krtObjectStore.Scope),
	}
}

func (h *Handler) mapKrtNetworkingToVersion(krtNetworking *krt.ProcessNetworking) *entity.ProcessNetworking {
	if krtNetworking == nil {
		return nil
	}

	return &entity.ProcessNetworking{
		TargetPort:      krtNetworking.TargetPort,
		DestinationPort: krtNetworking.DestinationPort,
		Protocol:        entity.NetworkingProtocol(string(krtNetworking.Protocol)),
	}
}

func (h *Handler) mapKrtResourceLimitToVersion(resourceLimit *krt.ResourceLimit) *entity.ResourceLimit {
	if resourceLimit == nil {
		return nil
	}

	return &entity.ResourceLimit{
		Request: resourceLimit.Request,
		Limit:   resourceLimit.Limit,
	}
}

func (h *Handler) mapKrtResourceLimitsToVersion(krtResourceLimits *krt.ProcessResourceLimits) *entity.ProcessResourceLimits {
	if krtResourceLimits == nil {
		return nil
	}

	return &entity.ProcessResourceLimits{
		CPU:    h.mapKrtResourceLimitToVersion(krtResourceLimits.CPU),
		Memory: h.mapKrtResourceLimitToVersion(krtResourceLimits.Memory),
	}
}
