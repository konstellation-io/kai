//go:build unit

package manager_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
)

func TestCreateKVStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).Return().AnyTimes()
	client := mocks.NewMockClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	const (
		testProductID         = "test-product"
		testVersionName       = "test-version"
		testWorkflowID        = "test-workflow"
		testProcessName       = "test-process"
		projectKeyValueStore  = "key-store_test-product_test-version"
		workflowKeyValueStore = "key-store_test-product_test-version_test-workflow"
		processKeyValueStore  = "key-store_test-product_test-version_test-workflow_test-process"
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
				fmt.Sprintf("key-store_%s_%s", testProductID, testVersionName),
				fmt.Sprintf("key-store_%s_%s_%s", testProductID, testVersionName, testWorkflowID),
				fmt.Sprintf("key-store_%s_%s_%s_%s", testProductID, testVersionName, testWorkflowID, testProcessName),
			},
			expectedWorkflowsKVCfg: &entity.VersionKeyValueStores{
				ProjectStore: projectKeyValueStore,
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
				fmt.Sprintf("key-store_%s_%s", testProductID, testVersionName),
				fmt.Sprintf("key-store_%s_%s_%s", testProductID, testVersionName, testWorkflowID),
			},
			expectedWorkflowsKVCfg: &entity.VersionKeyValueStores{
				ProjectStore: projectKeyValueStore,
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

			workflowsKVCfg, err := natsManager.CreateKeyValueStores(testProductID, testVersionName, tc.workflows)
			if tc.wantError {
				assert.ErrorIs(t, err, tc.wantedError)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedWorkflowsKVCfg, workflowsKVCfg)
		})
	}
}
