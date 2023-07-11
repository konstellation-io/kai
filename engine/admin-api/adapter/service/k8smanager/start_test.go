//go:build unit

package k8smanager_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/k8smanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

type startRequestMatcher struct {
	expectedStartRequest *versionpb.StartRequest
}

func newStartRequestMatcher(expectedStreamConfig *versionpb.StartRequest) *startRequestMatcher {
	return &startRequestMatcher{
		expectedStartRequest: expectedStreamConfig,
	}
}
func (m startRequestMatcher) String() string {
	return fmt.Sprintf("is equal to %v", m.expectedStartRequest)
}

func (m startRequestMatcher) Matches(actual interface{}) bool {
	actualCfg, ok := actual.(*versionpb.StartRequest)
	if !ok {
		return false
	}

	return reflect.DeepEqual(actualCfg, m.expectedStartRequest)
}

type StartVersionTestSuite struct {
	suite.Suite
	cfg              *config.Config
	logger           logging.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *k8smanager.K8sVersionClient
}

func TestStartVersionTestSuite(t *testing.T) {
	suite.Run(t, new(StartVersionTestSuite))
}

func (s *StartVersionTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	cfg := &config.Config{}
	logger := mocks.NewMockLogger(mockController)
	mocks.AddLoggerExpects(logger)
	service := mocks.NewMockVersionServiceClient(mockController)

	k8sVersionClient, err := k8smanager.NewK8sVersionClient(cfg, logger, service)
	s.Require().NoError(err)

	s.cfg = cfg
	s.logger = logger
	s.mockService = service
	s.k8sVersionClient = k8sVersionClient
}

func (s *StartVersionTestSuite) TestStartVersion() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
				Name:  "test-object-store",
				Scope: entity.ObjectStoreScopeWorkflow,
			}).
			WithNetworking(&entity.ProcessNetworking{
				TargetPort:      8080,
				DestinationPort: 8080,
				Protocol:        "TCP",
			}).
			WithConfig([]entity.ConfigurationVariable{
				{Key: "test-key", Value: "test-value"},
			}).
			Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(&version)

		workflowStreamCfg = versionConfig.StreamsConfig.Workflows[workflow.Name]
		processStreamCfg  = workflowStreamCfg.Processes[process.Name]

		workflowKVStoreCfg = versionConfig.KeyValueStoresConfig.Workflows[workflow.Name]
		processKVStore     = workflowKVStoreCfg.Processes[process.Name]
	)

	req := &versionpb.StartRequest{
		ProductId:     productID,
		VersionTag:    version.Version,
		KeyValueStore: versionConfig.KeyValueStoresConfig.KeyValueStore,
		Workflows: []*versionpb.Workflow{
			{
				Name:          workflow.Name,
				Stream:        versionConfig.StreamsConfig.Workflows[workflow.Name].Stream,
				KeyValueStore: versionConfig.KeyValueStoresConfig.Workflows[workflow.Name].KeyValueStore,
				Processes: []*versionpb.Process{
					{
						Name:          process.Name,
						Image:         process.Image,
						Gpu:           process.GPU,
						Subscriptions: processStreamCfg.Subscriptions,
						Subject:       processStreamCfg.Subject,
						Replicas:      process.Replicas,
						KeyValueStore: processKVStore,
						ObjectStore:   &process.ObjectStore.Name,
						Networking: &versionpb.Network{
							TargetPort: int32(process.Networking.TargetPort),
							Protocol:   process.Networking.Protocol,
							SourcePort: int32(process.Networking.DestinationPort),
						},
						Type: versionpb.ProcessType_ProcessTypeTask,
						Config: map[string]string{
							process.Config[0].Key: process.Config[0].Value,
						},
					},
				},
			},
		},
	}

	customMatcher := newStartRequestMatcher(req)
	s.mockService.EXPECT().Start(ctx, customMatcher).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Require().NoError(err)
}

func (s *StartVersionTestSuite) TestStartVersion_ClientError() {
	ctx := context.Background()

	var (
		productID     = "test-product"
		version       = testhelpers.NewVersionBuilder().Build()
		versionConfig = s.getConfigForVersion(&version)
	)

	expectedError := errors.New("client error")

	s.mockService.EXPECT().Start(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Assert().ErrorIs(err, expectedError)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_WorkflowStreamFound() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
				Name:  "test-object-store",
				Scope: entity.ObjectStoreScopeWorkflow,
			}).
			WithNetworking(&entity.ProcessNetworking{
				TargetPort:      8080,
				DestinationPort: 8080,
				Protocol:        "TCP",
			}).
			Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(&version)
	)

	// override default workflow to empty map
	versionConfig.StreamsConfig.Workflows = map[string]entity.WorkflowStreamConfig{}

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowStreamNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_NoWorkflowKeyValueStoreFound() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.KeyValueStoresConfig.Workflows = entity.WorkflowsKeyValueStoresConfig{}

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowKVStoreNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_NoWorkflowObjectStoreFound() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.ObjectStoresConfig.Workflows = map[string]entity.WorkflowObjectStoresConfig{}

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowObjectStoreNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_ProcessStreamNotFound() {
	ctx := context.Background()

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.StreamsConfig.Workflows[workflow.Name] = entity.WorkflowStreamConfig{
		Stream:    "stream",
		Processes: map[string]entity.ProcessStreamConfig{},
	}

	err := s.k8sVersionClient.Start(ctx, productID, &version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrProcessStreamNotFound)
}

func (s *StartVersionTestSuite) getConfigForVersion(version *entity.Version) *entity.VersionConfig {
	var (
		workflow = version.Workflows[0]
		process  = version.Workflows[0].Processes[0]
	)

	streamConfig := &entity.VersionStreamsConfig{
		Workflows: map[string]entity.WorkflowStreamConfig{
			workflow.Name: {
				Stream: "test-stream",
				Processes: map[string]entity.ProcessStreamConfig{
					process.Name: {
						Subject:       "process-subject",
						Subscriptions: []string{"another-node-subject"},
					},
				},
			},
		},
	}

	objectStoreConfig := &entity.VersionObjectStoresConfig{
		Workflows: map[string]entity.WorkflowObjectStoresConfig{
			workflow.Name: {
				Processes: map[string]string{
					process.Name: "test-object-store",
				},
			},
		},
	}

	keyValueStoresConfig := &entity.KeyValueStoresConfig{
		KeyValueStore: "test-product-kv-store",
		Workflows: map[string]*entity.WorkflowKeyValueStores{
			workflow.Name: {
				KeyValueStore: "test-workflow-kv-store",
				Processes: map[string]string{
					process.Name: "test-process-kv-store",
				},
			},
		},
	}

	return &entity.VersionConfig{
		StreamsConfig:        streamConfig,
		ObjectStoresConfig:   objectStoreConfig,
		KeyValueStoresConfig: keyValueStoresConfig,
	}
}
