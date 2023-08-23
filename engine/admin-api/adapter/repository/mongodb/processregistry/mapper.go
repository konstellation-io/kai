package processregistry

import (
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapDTOToEntity(dto *processRegistryDTO) *entity.ProcessRegistry {
	return &entity.ProcessRegistry{
		ID:         dto.ID,
		Name:       dto.Name,
		Version:    dto.Version,
		Type:       dto.Type,
		Image:      dto.Image,
		UploadDate: time.UnixMilli(dto.UploadDate).UTC(),
		Owner:      dto.Owner,
	}
}

func mapEntityToDTO(processEntity *entity.ProcessRegistry) *processRegistryDTO {
	return &processRegistryDTO{
		ID:         processEntity.ID,
		Name:       processEntity.Name,
		Version:    processEntity.Version,
		Type:       processEntity.Type,
		Image:      processEntity.Image,
		UploadDate: processEntity.UploadDate.UnixMilli(),
		Owner:      processEntity.Owner,
	}
}
