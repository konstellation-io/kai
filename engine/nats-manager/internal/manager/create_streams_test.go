//go:build unit

package manager_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/manager"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
)

func TestCreateStreams(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"
	testWorkflowName := "test-workflow"
	testStreamName := "test-product_test-version_test-workflow"
	testProcess := "test-process"

	testProcessSubject := fmt.Sprintf("%s.%s", testStreamName, testProcess)

	workflows := []entity.Workflow{
		NewWorkflowBuilder().
			WithID(testWorkflowName).
			WithProcessName(testProcess).
			Build(),
	}

	expectedWorkflowsStreamsCfg := entity.WorkflowsStreamsConfig{
		testWorkflowName: &entity.StreamConfig{
			Stream: testStreamName,
			Processes: entity.ProcessesStreamConfig{
				testProcess: entity.ProcessStreamConfig{
					Subject:       testProcessSubject,
					Subscriptions: []string{},
				},
			},
		},
	}

	customMatcher := newStreamConfigMatcher(expectedWorkflowsStreamsCfg[testWorkflowName])

	client.EXPECT().CreateStream(customMatcher).Return(nil)
	workflowsStreamsCfg, err := natsManager.CreateStreams(testProductID, testVersionName, workflows)
	assert.Nil(t, err)
	assert.Equal(t, expectedWorkflowsStreamsCfg, workflowsStreamsCfg)
}

func TestCreateStreams_ClientFails(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)
	client := mocks.NewMockClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	const (
		testProductID    = "test-product"
		testVersionName  = "test-version"
		testWorkflowName = "test-workflow"
		testProcess      = "test-process"
	)

	workflows := []entity.Workflow{
		NewWorkflowBuilder().
			WithID(testWorkflowName).
			WithProcessName(testProcess).
			Build(),
	}

	expectedError := fmt.Errorf("stream already exists")

	client.EXPECT().CreateStream(gomock.Any()).Return(fmt.Errorf("stream already exists"))
	workflowsStreamsConfig, err := natsManager.CreateStreams(testProductID, testVersionName, workflows)
	assert.Error(t, expectedError, err)
	assert.Nil(t, workflowsStreamsConfig)
}

func TestCreateStreams_FailsIfNoWorkflowsAreDefined(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mocks.NewMockLogger(ctrl)
	client := mocks.NewMockClient(ctrl)
	natsManager := manager.NewNatsManager(logger, client)

	testProductID := "test-product"
	testVersionName := "test-version"

	var workflows []entity.Workflow

	_, err := natsManager.CreateStreams(testProductID, testVersionName, workflows)
	assert.EqualError(t, err, "no workflows defined")
}
