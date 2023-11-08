//go:build unit

package manager_test

import (
	"testing"

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

	logger := mocks.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any()).Return().AnyTimes()
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
