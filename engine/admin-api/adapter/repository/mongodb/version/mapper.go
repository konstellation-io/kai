package version

import (
	"github.com/konstellation-io/krt/pkg/krt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapDTOToEntity(dto *versionDTO) *entity.Version {
	return &entity.Version{
		Krt: &krt.Krt{
			Name:        dto.Name,
			Description: dto.Description,
			Config:      dto.Config,
			Workflows:   mapDTOToEntityWorkflows(dto.Workflows),
		},
		ID: dto.ID,

		CreationDate:   dto.CreationDate,
		CreationAuthor: dto.CreationAuthor,

		PublicationDate:   dto.PublicationDate,
		PublicationAuthor: dto.PublicationAuthor,

		Status: dto.Status,
		Errors: dto.Errors,
	}
}

func mapDTOToEntityWorkflows(dtos []workflowDTO) []krt.Workflow {
	workflows := make([]krt.Workflow, len(dtos))

	for _, dto := range dtos {
		workflows = append(workflows, krt.Workflow{
			Name:      dto.ID,
			Type:      dto.Type,
			Config:    dto.Config,
			Processes: mapDTOToEntityProcesses(dto.Processes),
			Stream:    dto.Stream,
		})
	}

	return workflows
}

func mapDTOToEntityProcesses(dtos []processDTO) []krt.Process {
	processes := make([]krt.Process, len(dtos))

	for _, dto := range dtos {
		replicas := int(dto.Replicas)
		gpu := dto.GPU
		processes = append(processes, krt.Process{
			Name:          dto.ID,
			Type:          dto.Type,
			Image:         dto.Image,
			Replicas:      &replicas,
			GPU:           &gpu,
			Config:        dto.Config,
			ObjectStore:   mapDTOToEntityProcessObjectStore(dto.ObjectStore),
			Secrets:       dto.Secrets,
			Subscriptions: dto.Subscriptions,
			Networking:    mapDTOToEntityProcessNetworking(dto.Networking),
			Status:        dto.Status,
		})
	}

	return processes
}

func mapDTOToEntityProcessObjectStore(dto *processObjectStoreDTO) *krt.ProcessObjectStore {
	if dto == nil {
		return nil
	}

	return &krt.ProcessObjectStore{
		Name:  dto.Name,
		Scope: dto.Scope,
	}
}

func mapDTOToEntityProcessNetworking(dto *processNetworkingDTO) *krt.ProcessNetworking {
	if dto == nil {
		return nil
	}

	return &krt.ProcessNetworking{
		TargetPort:          dto.TargetPort,
		TargetProtocol:      dto.TargetProtocol,
		DestinationPort:     dto.DestinationPort,
		DestinationProtocol: dto.DestinationProtocol,
	}
}

func mapEntityToDTO(versionEntity *entity.Version) *versionDTO {
	return &versionDTO{
		ID:          versionEntity.ID,
		Name:        versionEntity.Name,
		Description: versionEntity.Description,
		Config:      versionEntity.Config,
		Workflows:   mapEntityToDTOWorkflows(versionEntity.Workflows),

		CreationDate:   versionEntity.CreationDate,
		CreationAuthor: versionEntity.CreationAuthor,

		PublicationDate:   versionEntity.PublicationDate,
		PublicationAuthor: versionEntity.PublicationAuthor,

		Status: versionEntity.Status,

		Errors: versionEntity.Errors,
	}
}

func mapEntityToDTOWorkflows(workflows []krt.Workflow) []workflowDTO {
	dtos := make([]workflowDTO, len(workflows))

	for _, workflow := range workflows {
		dtos = append(dtos, workflowDTO{
			ID:        workflow.Name,
			Type:      workflow.Type,
			Config:    workflow.Config,
			Processes: mapEntityToDTOProcesses(workflow.Processes),
			Stream:    workflow.Stream,
		})
	}

	return dtos
}

func mapEntityToDTOProcesses(processes []krt.Process) []processDTO {
	dtos := make([]processDTO, len(processes))

	for _, process := range processes {
		dtos = append(dtos, processDTO{
			ID:            process.Name,
			Type:          process.Type,
			Image:         process.Image,
			Replicas:      int32(*process.Replicas),
			GPU:           *process.GPU,
			Config:        process.Config,
			ObjectStore:   mapEntityToDTOProcessObjectStore(process.ObjectStore),
			Secrets:       process.Secrets,
			Subscriptions: process.Subscriptions,
			Networking:    mapEntityToDTOProcessNetworking(process.Networking),
		})
	}

	return dtos
}

func mapEntityToDTOProcessObjectStore(objectStore *krt.ProcessObjectStore) *processObjectStoreDTO {
	if objectStore == nil {
		return nil
	}

	return &processObjectStoreDTO{
		Name:  objectStore.Name,
		Scope: objectStore.Scope,
	}
}

func mapEntityToDTOProcessNetworking(networking *krt.ProcessNetworking) *processNetworkingDTO {
	if networking == nil {
		return nil
	}

	return &processNetworkingDTO{
		TargetPort:          networking.TargetPort,
		TargetProtocol:      networking.TargetProtocol,
		DestinationPort:     networking.DestinationPort,
		DestinationProtocol: networking.DestinationProtocol,
	}
}
