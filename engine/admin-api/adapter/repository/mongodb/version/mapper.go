package version

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapDTOToEntity(dto *versionDTO) *entity.Version {
	return &entity.Version{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Config:      mapDTOConfigToEntityConfig(dto.Config),
		Workflows:   mapDTOToEntityWorkflows(dto.Workflows),

		CreationDate:   dto.CreationDate,
		CreationAuthor: dto.CreationAuthor,

		PublicationDate:   dto.PublicationDate,
		PublicationAuthor: dto.PublicationAuthor,

		Status: entity.VersionStatus(dto.Status),
		Errors: dto.Errors,
	}
}

func mapDTOToEntityWorkflows(dtos []workflowDTO) []entity.Workflow {
	workflows := make([]entity.Workflow, len(dtos))

	for _, dto := range dtos {
		workflows = append(workflows, entity.Workflow{
			ID:        dto.ID,
			Name:      dto.Name,
			Type:      entity.WorkflowType(dto.Type),
			Config:    mapDTOConfigToEntityConfig(dto.Config),
			Processes: mapDTOToEntityProcesses(dto.Processes),
		})
	}

	return workflows
}

func mapDTOToEntityProcesses(dtos []processDTO) []entity.Process {
	processes := make([]entity.Process, len(dtos))

	for _, dto := range dtos {
		processes = append(processes, entity.Process{
			Name:          dto.ID,
			Type:          entity.ProcessType(dto.Type),
			Image:         dto.Image,
			Replicas:      dto.Replicas,
			GPU:           dto.GPU,
			Config:        mapDTOConfigToEntityConfig(dto.Config),
			ObjectStore:   mapDTOToEntityProcessObjectStore(dto.ObjectStore),
			Secrets:       mapDTOConfigToEntityConfig(dto.Secrets),
			Subscriptions: dto.Subscriptions,
			Networking:    mapDTOToEntityProcessNetworking(dto.Networking),
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
		Protocol:        dto.Protocol,
	}
}

func mapDTOConfigToEntityConfig(config []ConfigurationVariable) []entity.ConfigurationVariable {
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
		ID:          versionEntity.ID,
		Name:        versionEntity.Name,
		Description: versionEntity.Description,
		Config:      mapEntityConfigToDTOConfig(versionEntity.Config),
		Workflows:   mapEntityToDTOWorkflows(versionEntity.Workflows),

		CreationDate:   versionEntity.CreationDate,
		CreationAuthor: versionEntity.CreationAuthor,

		PublicationDate:   versionEntity.PublicationDate,
		PublicationAuthor: versionEntity.PublicationAuthor,

		Status: versionEntity.Status.String(),

		Errors: versionEntity.Errors,
	}
}

func mapEntityToDTOWorkflows(workflows []entity.Workflow) []workflowDTO {
	dtos := make([]workflowDTO, len(workflows))

	for _, workflow := range workflows {
		dtos = append(dtos, workflowDTO{
			ID:        workflow.Name,
			Type:      workflow.Type.String(),
			Config:    mapEntityConfigToDTOConfig(workflow.Config),
			Processes: mapEntityToDTOProcesses(workflow.Processes),
		})
	}

	return dtos
}

func mapEntityToDTOProcesses(processes []entity.Process) []processDTO {
	dtos := make([]processDTO, len(processes))

	for _, process := range processes {
		dtos = append(dtos, processDTO{
			ID:            process.Name,
			Type:          process.Type.String(),
			Image:         process.Image,
			Replicas:      process.Replicas,
			GPU:           process.GPU,
			Config:        mapEntityConfigToDTOConfig(process.Config),
			ObjectStore:   mapEntityToDTOProcessObjectStore(process.ObjectStore),
			Secrets:       mapEntityConfigToDTOConfig(process.Secrets),
			Subscriptions: process.Subscriptions,
			Networking:    mapEntityToDTOProcessNetworking(process.Networking),
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
		Protocol:        networking.Protocol,
	}
}

func mapEntityConfigToDTOConfig(config []entity.ConfigurationVariable) []ConfigurationVariable {
	dtoConfig := make([]ConfigurationVariable, len(config))

	for i, c := range config {
		dtoConfig[i] = ConfigurationVariable{
			Key:   c.Key,
			Value: c.Value,
		}
	}

	return dtoConfig
}
