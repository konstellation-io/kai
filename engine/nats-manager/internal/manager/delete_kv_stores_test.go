//go:build unit

package manager_test

import (
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/stretchr/testify/require"
)

func TestDeleteKeyValueStores(t *testing.T) {
	const (
		testProductID         = "test-product"
		testVersionTag        = "v1.0.0"
		testWorkflowID        = "test-workflow"
		testProcessName       = "test-process"
		versionKeyValueStore  = "key-store_test-product_v1_0_0"
		workflowKeyValueStore = "key-store_test-product_v1_0_0_test-workflow"
		processKeyValueStore  = "key-store_test-product_v1_0_0_test-workflow_test-process"
	)

	ctrl := gomock.NewController(t)

	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	workflows := []entity.Workflow{
		NewWorkflowBuilder().
			WithID(testWorkflowID).
			WithProcessName(testProcessName).
			Build(),
	}

	client.EXPECT().DeleteKeyValueStore(versionKeyValueStore).Return(nil)
	client.EXPECT().DeleteKeyValueStore(workflowKeyValueStore).Return(nil)
	client.EXPECT().DeleteKeyValueStore(processKeyValueStore).Return(nil)

	err := natsManager.DeleteVersionKeyValueStores(testProductID, testVersionTag, workflows)
	require.NoError(t, err)
}

func TestDeleteGlobalKeyValueStore(t *testing.T) {
	const (
		testProductID       = "test-product"
		globalKeyValueStore = "key-store_test-product"
	)

	ctrl := gomock.NewController(t)

	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	client.EXPECT().DeleteKeyValueStore(globalKeyValueStore).Return(nil)

	err := natsManager.DeleteGlobalKeyValueStore(testProductID)
	require.NoError(t, err)
}

func TestDeleteGlobalKeyValueStore_ClientError(t *testing.T) {
	const (
		testProductID       = "test-product"
		globalKeyValueStore = "key-store_test-product"
	)

	expectedErr := errors.New("client error")

	ctrl := gomock.NewController(t)

	logger := testr.NewWithOptions(t, testr.Options{Verbosity: -1})
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	client.EXPECT().DeleteKeyValueStore(globalKeyValueStore).Return(expectedErr)

	err := natsManager.DeleteGlobalKeyValueStore(testProductID)
	require.ErrorIs(t, err, expectedErr)
}
