//go:build unit

package processregistry

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

var domainProcessRegistry = &entity.ProcessRegistry{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	Image:      "test_image",
	UploadDate: testMapperUploadDate,
	Owner:      userID,
}

var DTOProcessRegistry = &processRegistryDTO{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	Image:      "test_image",
	UploadDate: testMapperUploadDate,
	Owner:      userID,
}

func TestMapDTOToEntity(t *testing.T) {
	obtainedDomainProcessRegistry := mapDTOToEntity(DTOProcessRegistry)
	assert.Equal(t, domainProcessRegistry, obtainedDomainProcessRegistry)
}

func TestMapEntityToDTO(t *testing.T) {
	obtainedDTOProcessRegistry := mapEntityToDTO(domainProcessRegistry)
	assert.Equal(t, DTOProcessRegistry, obtainedDTOProcessRegistry)
}
