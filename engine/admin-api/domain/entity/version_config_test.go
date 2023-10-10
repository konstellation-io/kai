package entity_test

import (
	"testing"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewVersionConfig(t *testing.T) {
	expected := &entity.VersionConfig{
		StreamsConfig: &entity.VersionStreamsConfig{
			Workflows: map[string]entity.WorkflowStreamConfig{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
		ObjectStoresConfig: &entity.VersionObjectStoresConfig{
			Workflows: map[string]entity.WorkflowObjectStoresConfig{},
		},
		KeyValueStoresConfig: &entity.KeyValueStores{
			KeyValueStore: "versionKVStore",
			Workflows: map[string]*entity.WorkflowKeyValueStores{
				"workflow-1": {
					KeyValueStore: "workflowKVStore",
				},
			},
		},
	}

	actual, err := entity.NewVersionConfig(expected.StreamsConfig, expected.ObjectStoresConfig, expected.KeyValueStoresConfig)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewVersionConfig_ErrorIfStramConfigNil(t *testing.T) {
	expected := &entity.VersionConfig{
		KeyValueStoresConfig: &entity.KeyValueStores{
			KeyValueStore: "versionKVStore",
			Workflows: map[string]*entity.WorkflowKeyValueStores{
				"workflow-1": {
					KeyValueStore: "workflowKVStore",
				},
			},
		},
	}

	_, err := entity.NewVersionConfig(nil, nil, expected.KeyValueStoresConfig)
	assert.ErrorIs(t, err, entity.ErrNilVersionStreamConfig)
}

func TestNewVersionConfig_ErrorIfKeyValueStoreConfigNil(t *testing.T) {
	expected := &entity.VersionConfig{
		StreamsConfig: &entity.VersionStreamsConfig{
			Workflows: map[string]entity.WorkflowStreamConfig{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
	}

	_, err := entity.NewVersionConfig(expected.StreamsConfig, nil, nil)
	assert.ErrorIs(t, err, entity.ErrNilKeyValueStoreConfig)
}
