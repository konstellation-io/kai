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
	testVersionTag := "v1.0.0"
	testStreamName := "test-product_v1.0.0_TestWorkflow"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionTag))

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return([]string{testStreamName}, nil)
	client.EXPECT().DeleteStream(testStreamName).Return(nil)
	err := natsManager.DeleteStreams(testProductID, testVersionTag)
	assert.Nil(t, err)
}

func TestDeleteStreams_ErrorDeletingStream(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionTag := "v1.0.0"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionTag))
	expectedError := errors.New("error getting streams")

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return(nil, expectedError)
	err := natsManager.DeleteStreams(testProductID, testVersionTag)
	assert.ErrorIs(t, err, expectedError)
}
func TestDeleteStreams_ErrorGettingStreamsNames(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockNatsClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionTag := "v1.0.0"
	testStreamName := "test-product_v1.0.0_TestWorkflow"
	expectedVersionStreamPattern := regexp.MustCompile(fmt.Sprintf("^%s_%s_.*", testProductID, testVersionTag))
	expectedError := errors.New("error deleting streams")

	client.EXPECT().GetStreamNames(expectedVersionStreamPattern).Return([]string{testStreamName}, nil)
	client.EXPECT().DeleteStream(testStreamName).Return(expectedError)

	err := natsManager.DeleteStreams(testProductID, testVersionTag)
	assert.ErrorIs(t, err, expectedError)
}
