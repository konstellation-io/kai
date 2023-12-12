package processrepository

import (
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

func mapDTOToEntity(dto *registeredProcessDTO) *entity.RegisteredProcess {
	return &entity.RegisteredProcess{
		ID:         dto.ID,
		Name:       dto.Name,
		Version:    dto.Version,
		Type:       entity.ProcessType(dto.Type),
		Image:      dto.Image,
		UploadDate: dto.UploadDate,
		Owner:      dto.Owner,
		Status:     dto.Status,
		Logs:       dto.Logs,
		IsPublic:   dto.IsPublic,
	}
}

func mapEntityToDTO(processEntity *entity.RegisteredProcess) *registeredProcessDTO {
	return &registeredProcessDTO{
		ID:         processEntity.ID,
		Name:       processEntity.Name,
		Version:    processEntity.Version,
		Type:       processEntity.Type.String(),
		Image:      processEntity.Image,
		UploadDate: processEntity.UploadDate,
		Owner:      processEntity.Owner,
		Status:     processEntity.Status,
		Logs:       processEntity.Logs,
		IsPublic:   processEntity.IsPublic,
	}
}
