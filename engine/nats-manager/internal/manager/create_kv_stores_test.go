//go:build unit

package manager_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateKVStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).Return().AnyTimes()
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	const (
		testProductID         = "test-product"
		testVersionTag        = "v1.0.0"
		testWorkflowID        = "test-workflow"
		testProcessName       = "test-process"
		testCleanedVersionTag = "v1_0_0"
		versionKeyValueStore  = "key-store_test-product_v1_0_0"
		workflowKeyValueStore = "key-store_test-product_v1_0_0_test-workflow"
		processKeyValueStore  = "key-store_test-product_v1_0_0_test-workflow_test-process"
	)

	tests := []struct {
		name                   string
		workflows              []entity.Workflow
		expectedKVStores       []string
		expectedWorkflowsKVCfg *entity.VersionKeyValueStores
		wantError              bool
		wantedError            error
		clientError            bool
	}{
		{
			name: "Key value stores for a workflow with a process",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowID).
					WithProcessName(testProcessName).
					Build(),
			},
			expectedKVStores: []string{
				fmt.Sprintf("key-store_%s_%s", testProductID, testCleanedVersionTag),
				fmt.Sprintf("key-store_%s_%s_%s", testProductID, testCleanedVersionTag, testWorkflowID),
				fmt.Sprintf("key-store_%s_%s_%s_%s", testProductID, testCleanedVersionTag, testWorkflowID, testProcessName),
			},
			expectedWorkflowsKVCfg: &entity.VersionKeyValueStores{
				ProjectStore: versionKeyValueStore,
				WorkflowsStores: map[string]*entity.WorkflowKeyValueStores{
					testWorkflowID: {
						WorkflowStore: workflowKeyValueStore,
						Processes: map[string]string{
							testProcessName: processKeyValueStore,
						},
					},
				},
			},
			wantError:   false,
			wantedError: nil,
		},
		{
			name: "Key value stores for a workflow without a process",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowID).
					WithProcesses(nil).
					Build(),
			},
			expectedKVStores: []string{
				fmt.Sprintf("key-store_%s_%s", testProductID, testCleanedVersionTag),
				fmt.Sprintf("key-store_%s_%s_%s", testProductID, testCleanedVersionTag, testWorkflowID),
			},
			expectedWorkflowsKVCfg: &entity.VersionKeyValueStores{
				ProjectStore: versionKeyValueStore,
				WorkflowsStores: map[string]*entity.WorkflowKeyValueStores{
					testWorkflowID: {
						WorkflowStore: workflowKeyValueStore,
						Processes:     map[string]string{},
					},
				},
			},
			wantError:   false,
			wantedError: nil,
		},
		{
			name:                   "Key value stores without a workflow",
			workflows:              nil,
			expectedKVStores:       []string{},
			expectedWorkflowsKVCfg: nil,
			wantError:              true,
			wantedError:            internal.ErrNoWorkflowsDefined,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, expectedKVStore := range tc.expectedKVStores {
				client.EXPECT().CreateKeyValueStore(expectedKVStore).Return(nil)
			}

			workflowsKVCfg, err := natsManager.CreateVersionKeyValueStores(testProductID, testVersionTag, tc.workflows)
			if tc.wantError {
				assert.ErrorIs(t, err, tc.wantedError)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedWorkflowsKVCfg, workflowsKVCfg)
		})
	}
}

func TestCreateGlobalKVStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Return().AnyTimes()
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	const (
		testProductID       = "test-product"
		globalKeyValueStore = "key-store_test-product"
	)

	client.EXPECT().CreateKeyValueStore(globalKeyValueStore).Return(nil)

	globalKVStore, err := natsManager.CreateGlobalKeyValueStore(testProductID)
	require.NoError(t, err)

	assert.Equal(t, globalKeyValueStore, globalKVStore)
}

func TestCreateGlobalKVStore_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	logger.EXPECT().Infof(gomock.Any(), gomock.Any()).Return().AnyTimes()
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	const (
		testProductID       = "test-product"
		globalKeyValueStore = "key-store_test-product"
	)

	client.EXPECT().CreateKeyValueStore(globalKeyValueStore).Return(errors.New("mock error"))

	_, err := natsManager.CreateGlobalKeyValueStore(testProductID)
	assert.Error(t, err)
}
