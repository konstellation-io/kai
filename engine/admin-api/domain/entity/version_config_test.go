//go:build unit

package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVersionConfig(t *testing.T) {
	expected := &entity.VersionStreamingResources{
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamConfig{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
		ObjectStores: &entity.VersionObjectStores{
			Workflows: map[string]entity.WorkflowObjectStoresConfig{},
		},
		KeyValueStores: &entity.KeyValueStores{
			KeyValueStore: "versionKVStore",
			Workflows: map[string]*entity.WorkflowKeyValueStores{
				"workflow-1": {
					KeyValueStore: "workflowKVStore",
				},
			},
		},
	}

	actual, err := entity.NewVersionConfig(expected.Streams, expected.ObjectStores, expected.KeyValueStores)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewVersionConfig_ErrorIfStramConfigNil(t *testing.T) {
	expected := &entity.VersionStreamingResources{
		KeyValueStores: &entity.KeyValueStores{
			KeyValueStore: "versionKVStore",
			Workflows: map[string]*entity.WorkflowKeyValueStores{
				"workflow-1": {
					KeyValueStore: "workflowKVStore",
				},
			},
		},
	}

	_, err := entity.NewVersionConfig(nil, nil, expected.KeyValueStores)
	assert.ErrorIs(t, err, entity.ErrNilVersionStreamConfig)
}

func TestNewVersionConfig_ErrorIfKeyValueStoreConfigNil(t *testing.T) {
	expected := &entity.VersionStreamingResources{
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamConfig{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
	}

	_, err := entity.NewVersionConfig(expected.Streams, nil, nil)
	assert.ErrorIs(t, err, entity.ErrNilKeyValueStoreConfig)
}
