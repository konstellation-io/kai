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
	client := mocks.NewMockClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"
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
			expectedObjectStores: []string{fmt.Sprintf("%s_%s_%s", testProductID, testVersionName, testObjectStore)},
		},
		{
			name: "Project with multiple object stores",
			expectedObjectStores: []string{
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionName, testWorkflowName, testObjectStore),
				fmt.Sprintf("%s_%s_%s_another-object-store", testProductID, testVersionName, testWorkflowName),
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
				fmt.Sprintf("%s_%s_%s_%s", testProductID, testVersionName, testWorkflowName, testObjectStore),
			},
			deleteObjectStoreError: errors.New("error deleting object store"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filter := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionName))

			client.EXPECT().GetObjectStoreNames(filter).Return(tc.expectedObjectStores, tc.getObjectStoresError)
			for _, expectedObjStore := range tc.expectedObjectStores {
				client.EXPECT().DeleteObjectStore(expectedObjStore).Return(tc.deleteObjectStoreError)
			}
			err := natsManager.DeleteObjectStores(testProductID, testVersionName)
			if tc.getObjectStoresError != nil {
				assert.ErrorIs(t, err, tc.getObjectStoresError)
			} else {
				assert.ErrorIs(t, err, tc.deleteObjectStoreError)
			}
		})
	}
}
