//go:build unit

package manager_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteObjectStore(t *testing.T) {
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
		name                   string
		expectedObjectStores   []string
		getObjectStoresError   error
		deleteObjectStoreError error
	}{
		{
			name:                 "Project with 1 object store",
			expectedObjectStores: []string{fmt.Sprintf("%s_%s_%s", testProductID, testVersionTag, testObjectStore)},
		},
		{
			name: "Project with multiple object stores",
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
				fmt.Sprintf("%s_%s_%s_another-object-store", testProductID, testVersionTag, testWorkflowName),
			},
		},
		{
			name:                 "Error getting object store",
			expectedObjectStores: nil,
			getObjectStoresError: errors.New("error getting object stores"),
		},
		{
			name: "Error deleting object store",
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionTag, testWorkflowName, testObjectStore),
			},
			deleteObjectStoreError: errors.New("error deleting object store"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filter := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionTag))

			client.EXPECT().GetObjectStoreNames(filter).Return(tc.expectedObjectStores, tc.getObjectStoresError)
			for _, expectedObjStore := range tc.expectedObjectStores {
				client.EXPECT().DeleteObjectStore(expectedObjStore).Return(tc.deleteObjectStoreError)
			}
			err := natsManager.DeleteObjectStores(testProductID, testVersionTag)
			if tc.getObjectStoresError != nil {
				assert.ErrorIs(t, err, tc.getObjectStoresError)
			} else {
				assert.ErrorIs(t, err, tc.deleteObjectStoreError)
			}
		})
	}
}
