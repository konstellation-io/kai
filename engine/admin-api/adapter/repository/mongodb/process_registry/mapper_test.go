//go:build unit

package process_registry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

var (
	userID     = "admin"
	uploadDate = time.Now().Add(-time.Hour)
)

var domainProcessRegistry = &entity.ProcessRegistry{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	UploadDate: uploadDate,
	Owner:      userID,
}

var DTOProcessRegistry = &processRegistryDTO{
	ID:         "id_process",
	Name:       "test_trigger",
	Version:    "1.0.0",
	Type:       "trigger",
	UploadDate: uploadDate,
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
