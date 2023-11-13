//go:build unit

package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestProduct_HasVersionPublished(t *testing.T) {
	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(testhelpers.StrPointer("test-version")).
		Build()

	assert.True(t, product.HasVersionPublished())
}

func TestProduct_HasVersionPublished_NoVersionPublished(t *testing.T) {
	product := testhelpers.NewProductBuilder().
		Build()

	assert.False(t, product.HasVersionPublished())
}

func TestProduct_HasVersionPublished_NoVersionPublished_EmtpyString(t *testing.T) {
	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(nil).
		Build()

	assert.False(t, product.HasVersionPublished())
}

func TestProduct_UpdatePublishedVersion(t *testing.T) {
	publishedVersion := "test-version"
	product := testhelpers.NewProductBuilder().Build()

	product.UpdatePublishedVersion(publishedVersion)

	assert.Equal(t, publishedVersion, *product.PublishedVersion)
}

func TestProduct_RemovePublishedVersion(t *testing.T) {
	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(testhelpers.StrPointer("test-version")).
		Build()

	product.RemovePublishedVersion()

	assert.Nil(t, product.PublishedVersion)
}
