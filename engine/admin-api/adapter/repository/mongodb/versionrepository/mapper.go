package versionrepository

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapDTOToEntity(dto *versionDTO) *entity.Version {
	return &entity.Version{
		Tag:         dto.Tag,
		Description: dto.Description,
		Config:      mapDTOConfigToEntityConfig(dto.Config),
		Workflows:   mapDTOToEntityWorkflows(dto.Workflows),

		CreationDate:   dto.CreationDate,
		CreationAuthor: dto.CreationAuthor,

		PublicationDate:   dto.PublicationDate,
		PublicationAuthor: dto.PublicationAuthor,

		Status: entity.VersionStatus(dto.Status),
		Error:  dto.Error,
	}
}

func mapDTOToEntityWorkflows(dtos []workflowDTO) []entity.Workflow {
	workflows := make([]entity.Workflow, 0, len(dtos))

	for _, dto := range dtos {
		workflows = append(workflows, entity.Workflow{
			Name:      dto.Name,
			Type:      entity.WorkflowType(dto.Type),
			Config:    mapDTOConfigToEntityConfig(dto.Config),
			Processes: mapDTOToEntityProcesses(dto.Processes),
		})
	}

	return workflows
}

func mapDTOToEntityProcesses(dtos []processDTO) []entity.Process {
	processes := make([]entity.Process, 0, len(dtos))

	for _, dto := range dtos {
		processes = append(processes, entity.Process{
			Name:           dto.Name,
			Type:           entity.ProcessType(dto.Type),
			Image:          dto.Image,
			Replicas:       dto.Replicas,
			GPU:            dto.GPU,
			Config:         mapDTOConfigToEntityConfig(dto.Config),
			ObjectStore:    mapDTOToEntityProcessObjectStore(dto.ObjectStore),
			Secrets:        mapDTOConfigToEntityConfig(dto.Secrets),
			Subscriptions:  dto.Subscriptions,
			Networking:     mapDTOToEntityProcessNetworking(dto.Networking),
			ResourceLimits: mapDTOToEntityProcessResourceLimits(dto.ResourceLimits),
			NodeSelectors:  dto.NodeSelectors,
		})
	}

	return processes
}

func mapDTOToEntityProcessObjectStore(dto *processObjectStoreDTO) *entity.ProcessObjectStore {
	if dto == nil {
		return nil
	}

	return &entity.ProcessObjectStore{
		Name:  dto.Name,
		Scope: entity.ObjectStoreScope(dto.Scope),
	}
}

func mapDTOToEntityProcessNetworking(dto *processNetworkingDTO) *entity.ProcessNetworking {
	if dto == nil {
		return nil
	}

	return &entity.ProcessNetworking{
		TargetPort:      dto.TargetPort,
		DestinationPort: dto.DestinationPort,
		Protocol:        entity.NetworkingProtocol(dto.Protocol),
	}
}

func mapDTOToEntityProcessResourceLimits(dto *processResourceLimitsDTO) *entity.ProcessResourceLimits {
	if dto == nil {
		return nil
	}

	return &entity.ProcessResourceLimits{
		CPU:    mapDTOToEntityResourceLimit(dto.CPU),
		Memory: mapDTOToEntityResourceLimit(dto.Memory),
	}
}

func mapDTOToEntityResourceLimit(dto *resourceLimitDTO) *entity.ResourceLimit {
	if dto == nil {
		return nil
	}

	return &entity.ResourceLimit{
		Request: dto.Request,
		Limit:   dto.Limit,
	}
}

func mapDTOConfigToEntityConfig(config []configurationVariableDTO) []entity.ConfigurationVariable {
	if config == nil {
		return nil
	}

	entityConfig := make([]entity.ConfigurationVariable, len(config))

	for i, c := range config {
		entityConfig[i] = entity.ConfigurationVariable{
			Key:   c.Key,
			Value: c.Value,
		}
	}

	return entityConfig
}

func mapEntityToDTO(versionEntity *entity.Version) *versionDTO {
	return &versionDTO{
		Tag:         versionEntity.Tag,
		Description: versionEntity.Description,
		Config:      mapEntityConfigToDTOConfig(versionEntity.Config),
		Workflows:   mapEntityToDTOWorkflows(versionEntity.Workflows),

		CreationDate:   versionEntity.CreationDate,
		CreationAuthor: versionEntity.CreationAuthor,

		PublicationDate:   versionEntity.PublicationDate,
		PublicationAuthor: versionEntity.PublicationAuthor,

		Status: versionEntity.Status.String(),

		Error: versionEntity.Error,
	}
}

func mapEntityToDTOWorkflows(workflows []entity.Workflow) []workflowDTO {
	dtos := make([]workflowDTO, len(workflows))
	idx := 0

	for _, workflow := range workflows {
		dtos[idx] = workflowDTO{
			Name:      workflow.Name,
			Type:      workflow.Type.String(),
			Config:    mapEntityConfigToDTOConfig(workflow.Config),
			Processes: mapEntityToDTOProcesses(workflow.Processes),
		}
		idx++
	}

	return dtos
}

func mapEntityToDTOProcesses(processes []entity.Process) []processDTO {
	dtos := make([]processDTO, 0, len(processes))

	for _, process := range processes {
		dtos = append(dtos, processDTO{
			Name:           process.Name,
			Type:           process.Type.String(),
			Image:          process.Image,
			Replicas:       process.Replicas,
			GPU:            process.GPU,
			Config:         mapEntityConfigToDTOConfig(process.Config),
			ObjectStore:    mapEntityToDTOProcessObjectStore(process.ObjectStore),
			Secrets:        mapEntityConfigToDTOConfig(process.Secrets),
			Subscriptions:  process.Subscriptions,
			Networking:     mapEntityToDTOProcessNetworking(process.Networking),
			ResourceLimits: mapEntityToDTOProcessResourceLimits(process.ResourceLimits),
			NodeSelectors:  process.NodeSelectors,
		})
	}

	return dtos
}

func mapEntityToDTOProcessObjectStore(objectStore *entity.ProcessObjectStore) *processObjectStoreDTO {
	if objectStore == nil {
		return nil
	}

	return &processObjectStoreDTO{
		Name:  objectStore.Name,
		Scope: objectStore.Scope.String(),
	}
}

func mapEntityToDTOProcessNetworking(networking *entity.ProcessNetworking) *processNetworkingDTO {
	if networking == nil {
		return nil
	}

	return &processNetworkingDTO{
		TargetPort:      networking.TargetPort,
		DestinationPort: networking.DestinationPort,
		Protocol:        string(networking.Protocol),
	}
}

func mapEntityToDTOProcessResourceLimits(resouceLimits *entity.ProcessResourceLimits) *processResourceLimitsDTO {
	if resouceLimits == nil {
		return nil
	}

	return &processResourceLimitsDTO{
		CPU:    mapEntityToDTOresourceLimit(resouceLimits.CPU),
		Memory: mapEntityToDTOresourceLimit(resouceLimits.Memory),
	}
}

func mapEntityToDTOresourceLimit(cpu *entity.ResourceLimit) *resourceLimitDTO {
	if cpu == nil {
		return nil
	}

	return &resourceLimitDTO{
		Request: cpu.Request,
		Limit:   cpu.Limit,
	}
}

func mapEntityConfigToDTOConfig(config []entity.ConfigurationVariable) []configurationVariableDTO {
	if config == nil {
		return nil
	}

	dtoConfig := make([]configurationVariableDTO, len(config))

	for i, c := range config {
		dtoConfig[i] = configurationVariableDTO{
			Key:   c.Key,
			Value: c.Value,
		}
	}

	return dtoConfig
}
