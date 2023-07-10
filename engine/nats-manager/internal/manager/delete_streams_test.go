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

func TestDeleteStreams(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"
	testStreamName := "test-product_test-version_TestWorkflow"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionName))

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return([]string{testStreamName}, nil)
	client.EXPECT().DeleteStream(testStreamName).Return(nil)
	err := natsManager.DeleteStreams(testProductID, testVersionName)
	assert.Nil(t, err)
}

func TestDeleteStreams_ErrorDeletingStream(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionName))
	expectedError := errors.New("error getting streams")

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return(nil, expectedError)
	err := natsManager.DeleteStreams(testProductID, testVersionName)
	assert.ErrorIs(t, err, expectedError)
}
func TestDeleteStreams_ErrorGettingStreamsNames(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"
	testStreamName := "test-product_test-version_TestWorkflow"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionName))
	expectedError := errors.New("error deleting streams")

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return([]string{testStreamName}, nil)
	client.EXPECT().DeleteStream(testStreamName).Return(expectedError)

	err := natsManager.DeleteStreams(testProductID, testVersionName)
	assert.ErrorIs(t, err, expectedError)
}
