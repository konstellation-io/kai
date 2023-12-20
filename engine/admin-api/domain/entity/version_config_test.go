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
			Workflows: map[string]entity.WorkflowStreamResources{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
		ObjectStores: &entity.VersionObjectStores{
			Workflows: map[string]entity.WorkflowObjectStoresConfig{},
		},
		KeyValueStores: &entity.KeyValueStores{
			VersionKeyValueStore: "versionKVStore",
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
			VersionKeyValueStore: "versionKVStore",
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
			Workflows: map[string]entity.WorkflowStreamResources{
				"workflow-1": {
					Stream: "stream-1",
				},
			},
		},
	}

	_, err := entity.NewVersionConfig(expected.Streams, nil, nil)
	assert.ErrorIs(t, err, entity.ErrNilKeyValueStoreConfig)
}

func TestGetWorkflowStreamConfig(t *testing.T) {
	var (
		workflowName = "workflow-1"

		expectedWorkflowStreamResources = &entity.WorkflowStreamResources{
			Stream: "stream-1",
		}
	)
	versionResources := &entity.VersionStreamingResources{
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamResources{
				workflowName: *expectedWorkflowStreamResources,
			},
		},
	}

	workflowStream, err := versionResources.GetWorkflowStream(workflowName)
	require.NoError(t, err)
	assert.Equal(t, expectedWorkflowStreamResources, workflowStream)
}

func TestGetWorkflowStreamConfig_WorkflowNotFound(t *testing.T) {
	versionResources := &entity.VersionStreamingResources{
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamResources{},
		},
	}

	_, err := versionResources.GetWorkflowStream("unknown-workflow")
	assert.ErrorIs(t, err, entity.ErrWorkflowStreamNotFound)
}

func TestGetKVStoreConfig(t *testing.T) {
	var (
		workflowName = "workflow-1"

		expectedWorkflowStreamResources = &entity.WorkflowStreamResources{
			Stream: "stream-1",
		}
	)
	versionResources := &entity.VersionStreamingResources{
		Streams: &entity.VersionStreams{
			Workflows: map[string]entity.WorkflowStreamResources{
				workflowName: *expectedWorkflowStreamResources,
			},
		},
	}

	workflowStream, err := versionResources.GetWorkflowStream(workflowName)
	require.NoError(t, err)
	assert.Equal(t, expectedWorkflowStreamResources, workflowStream)
}
