//go:build unit

package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestVersion_SetPublishStatus(t *testing.T) {
	publicationAuthor := "test-user"
	version := testhelpers.NewVersionBuilder().Build()

	version.SetPublishStatus(publicationAuthor)

	assert.Equal(t, entity.VersionStatusPublished, version.Status)
	assert.Equal(t, publicationAuthor, *version.PublicationAuthor)
}

func TestVersion_UnsetPublishStatus(t *testing.T) {
	version := testhelpers.NewVersionBuilder().Build()

	version.UnsetPublishStatus()

	assert.Equal(t, entity.VersionStatusStarted, version.Status)
	assert.Nil(t, version.PublicationAuthor)
}
