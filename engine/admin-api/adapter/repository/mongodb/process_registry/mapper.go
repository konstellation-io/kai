package process_registry

import "github.com/konstellation-io/kai/engine/admin-api/domain/entity"

func mapDTOToEntity(dto *processRegistryDTO) *entity.ProcessRegistry {
	return &entity.ProcessRegistry{
		ID:         dto.ID,
		Name:       dto.Name,
		Version:    dto.Version,
		Type:       dto.Type,
		UploadDate: dto.UploadDate,
		Owner:      dto.Owner,
	}
}

func mapEntityToDTO(entity *entity.ProcessRegistry) *processRegistryDTO {
	return &processRegistryDTO{
		ID:         entity.ID,
		Name:       entity.Name,
		Version:    entity.Version,
		Type:       entity.Type,
		UploadDate: entity.UploadDate,
		Owner:      entity.Owner,
	}
}
