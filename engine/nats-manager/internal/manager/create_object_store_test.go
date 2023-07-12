//go:build unit

package manager_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateObjectStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionTag := "v1.0.0"
	testWorkflowName := "test-workflow"
	testObjectStore := "test-object-store"

	tests := []struct {
		name                 string
		workflows            []entity.Workflow
		expectedObjectStores []string
		expectedError        error
	}{
		{
			name: "Object store with project scope",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeProject,
						},
					).
					Build(),
			},
			expectedObjectStores: []string{fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore)},
		},
		{
			name: "Object store with workflow scope",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeWorkflow,
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
			},
		},
		{
			name: "Invalid object store name",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  "",
							Scope: entity.ObjStoreScopeWorkflow,
						},
					).
					Build(),
			},
			expectedObjectStores: nil,
			expectedError:        internal.ErrInvalidObjectStoreName,
		},
		{
			name: "Invalid object store scope",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: "invalid",
						},
					).
					Build(),
			},
			expectedObjectStores: nil,
			expectedError:        internal.ErrInvalidObjectStoreScope,
		},
		{
			name: "Process without object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					Build(),
			},
			expectedObjectStores: nil,
		},
		{
			name: "Multiple workflows with different workflow scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeWorkflow,
						},
					).
					Build(),
				NewWorkflowBuilder().
					WithID("another-workflow").
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeWorkflow,
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
				fmt.Sprintf("%s_%s_another-workflow_%s", testProductID, testVersionTag, testObjectStore),
			},
		},
		{
			name: "Multiple workflows with the same project scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeProject,
						},
					).
					Build(),
				NewWorkflowBuilder().
					WithID("another-workflow").
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeProject,
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
			},
		},
		{
			name: "Multiple workflows with different project scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  testObjectStore,
							Scope: entity.ObjStoreScopeProject,
						},
					).
					Build(),
				NewWorkflowBuilder().
					WithID("another-workflow").
					WithProcessObjectStore(
						&entity.ObjectStore{
							Name:  "another-object-store",
							Scope: entity.ObjStoreScopeProject,
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
				fmt.Sprintf("%s_%s_another-object-store", testProductID, testVersionTag),
			},
		},
		{
			name: "Multiple processes in workflow with same workflow scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcesses(
						[]entity.Process{
							{
								Name: "test-process-1",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeWorkflow,
								},
							},
							{
								Name: "test-process-2",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeWorkflow,
								},
							},
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
			},
		},
		{
			name: "Multiple processes in workflow with different workflow scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcesses(
						[]entity.Process{
							{
								Name: "test-process-1",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeWorkflow,
								},
							},
							{
								Name: "test-process-2",
								ObjectStore: &entity.ObjectStore{
									Name:  "another-object-store",
									Scope: entity.ObjStoreScopeWorkflow,
								},
							},
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
				fmt.Sprintf("%s_%s_%s_another-object-store", testProductID, testVersionTag, testWorkflowName),
			},
		},
		{
			name: "Multiple processes in workflow with same project scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcesses(
						[]entity.Process{
							{
								Name: "test-process-1",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeProject,
								},
							},
							{
								Name: "test-process-2",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeProject,
								},
							},
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
			},
		},
		{
			name: "Multiple processes in workflow with different project scoped object store",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithID(testWorkflowName).
					WithProcesses(
						[]entity.Process{
							{
								Name: "test-process-1",
								ObjectStore: &entity.ObjectStore{
									Name:  testObjectStore,
									Scope: entity.ObjStoreScopeProject,
								},
							},
							{
								Name: "test-process-2",
								ObjectStore: &entity.ObjectStore{
									Name:  "another-object-store",
									Scope: entity.ObjStoreScopeProject,
								},
							},
						},
					).
					Build(),
			},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore),
				fmt.Sprintf("%s_%s_another-object-store", testProductID, testVersionTag),
			},
		},
		{
			name: "nats client error",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithProcessObjectStore(&entity.ObjectStore{
						Name:  testObjectStore,
						Scope: entity.ObjStoreScopeWorkflow,
					}).Build()},
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
			},
			expectedError: fmt.Errorf("nats client error"),
		},
		{
			name: "invalid scope error",
			workflows: []entity.Workflow{
				NewWorkflowBuilder().
					WithProcessObjectStore(&entity.ObjectStore{
						Name:  testObjectStore,
						Scope: entity.ObjectStoreScope("invalid"),
					}).Build()},
			expectedError: internal.ErrInvalidObjectStoreScope,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, expectedObjStore := range tc.expectedObjectStores {
				client.EXPECT().CreateObjectStore(expectedObjStore).Return(tc.expectedError)
			}

			_, err := natsManager.CreateObjectStores(testProductID, testVersionTag, tc.workflows)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
