//go:build unit

package processrepository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

var (
	userID               = "admin"
	testMapperUploadDate = time.Now().Add(-time.Hour).Truncate(time.Millisecond).UTC()
)

var domainRegisteredProcess = &entity.RegisteredProcess{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	Image:      "test_image",
	UploadDate: testMapperUploadDate,
	Owner:      userID,
}

var DTORegisteredProcess = &registeredProcessDTO{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	Image:      "test_image",
	UploadDate: testMapperUploadDate,
	Owner:      userID,
}

func TestMapDTOToEntity(t *testing.T) {
	obtainedDomainRegisteredProcess := mapDTOToEntity(DTORegisteredProcess)
	assert.Equal(t, domainRegisteredProcess, obtainedDomainRegisteredProcess)
}

func TestMapEntityToDTO(t *testing.T) {
	obtainedDTORegisteredProcess := mapEntityToDTO(domainRegisteredProcess)
	assert.Equal(t, DTORegisteredProcess, obtainedDTORegisteredProcess)
}
