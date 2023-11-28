//go:build unit

package versionservice_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
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
	k8sVersionClient *versionservice.K8sVersionService
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

	k8sVersionClient, err := versionservice.New(cfg, logger, service)
	s.Require().NoError(err)

	s.cfg = cfg
	s.logger = logger
	s.mockService = service
	s.k8sVersionClient = k8sVersionClient
}

func (s *StartVersionTestSuite) TestStartVersion() {
	ctx := context.Background()

	var (
		product = testhelpers.NewProductBuilder().Build()

		process = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
				Name:  "test-object-store",
				Scope: entity.ObjectStoreScopeWorkflow,
			}).
			WithNetworking(&entity.ProcessNetworking{
				TargetPort:      8080,
				DestinationPort: 8080,
				Protocol:        "GRPC",
			}).
			WithResourceLimits(&entity.ProcessResourceLimits{
				CPU: &entity.ResourceLimit{
					Request: "100m",
					Limit:   "200m",
				},
				Memory: &entity.ResourceLimit{
					Request: "100Mi",
					Limit:   "200Mi",
				},
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

		versionConfig = s.getConfigForVersion(version)

		workflowStreamCfg = versionConfig.Streams.Workflows[workflow.Name]
		processStreamCfg  = workflowStreamCfg.Processes[process.Name]

		workflowKVStoreCfg = versionConfig.KeyValueStores.Workflows[workflow.Name]
		processKVStore     = workflowKVStoreCfg.Processes[process.Name]
	)

	req := &versionpb.StartRequest{
		ProductId:            product.ID,
		VersionTag:           version.Tag,
		GlobalKeyValueStore:  versionConfig.KeyValueStores.GlobalKeyValueStore,
		VersionKeyValueStore: versionConfig.KeyValueStores.VersionKeyValueStore,
		Workflows: []*versionpb.Workflow{
			{
				Name:          workflow.Name,
				Stream:        versionConfig.Streams.Workflows[workflow.Name].Stream,
				KeyValueStore: versionConfig.KeyValueStores.Workflows[workflow.Name].KeyValueStore,
				Type:          workflow.Type.String(),
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
							Protocol:   string(process.Networking.Protocol),
							SourcePort: int32(process.Networking.DestinationPort),
						},
						ResourceLimits: &versionpb.ProcessResourceLimits{
							Cpu: &versionpb.ResourceLimit{
								Request: process.ResourceLimits.CPU.Request,
								Limit:   process.ResourceLimits.CPU.Limit,
							},
							Memory: &versionpb.ResourceLimit{
								Request: process.ResourceLimits.Memory.Request,
								Limit:   process.ResourceLimits.Memory.Limit,
							},
						},
						Type: versionpb.ProcessType_ProcessTypeTask,
						Config: map[string]string{
							process.Config[0].Key: process.Config[0].Value,
						},
					},
				},
			},
		},
		MinioConfiguration: &versionpb.MinioConfiguration{
			Bucket: product.MinioConfiguration.Bucket,
		},
		ServiceAccount: &versionpb.ServiceAccount{
			Username: product.ServiceAccount.Username,
			Password: product.ServiceAccount.Password,
		},
	}

	customMatcher := newStartRequestMatcher(req)
	s.mockService.EXPECT().Start(ctx, customMatcher).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Require().NoError(err)
}

func (s *StartVersionTestSuite) TestStartVersion_ClientError() {
	ctx := context.Background()

	var (
		product       = testhelpers.NewProductBuilder().Build()
		version       = testhelpers.NewVersionBuilder().Build()
		versionConfig = s.getConfigForVersion(version)
	)

	expectedError := errors.New("client error")

	s.mockService.EXPECT().Start(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Assert().ErrorIs(err, expectedError)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_WorkflowStreamFound() {
	ctx := context.Background()

	var (
		product = testhelpers.NewProductBuilder().Build()

		process = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
				Name:  "test-object-store",
				Scope: entity.ObjectStoreScopeWorkflow,
			}).
			WithNetworking(&entity.ProcessNetworking{
				TargetPort:      8080,
				DestinationPort: 8080,
				Protocol:        "GRPC",
			}).
			WithResourceLimits(&entity.ProcessResourceLimits{
				CPU: &entity.ResourceLimit{
					Request: "100m",
					Limit:   "200m",
				},
				Memory: &entity.ResourceLimit{
					Request: "100Mi",
					Limit:   "200Mi",
				},
			}).
			Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(version)
	)

	// override default workflow to empty map
	versionConfig.Streams.Workflows = map[string]entity.WorkflowStreamResources{}

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowStreamNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_NoWorkflowKeyValueStoreFound() {
	ctx := context.Background()

	var (
		product = testhelpers.NewProductBuilder().Build()

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(version)
	)

	// override default workflow config to empty map
	versionConfig.KeyValueStores.Workflows = map[string]*entity.WorkflowKeyValueStores{}

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowKVStoreNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_NoWorkflowObjectStoreFound() {
	ctx := context.Background()

	var (
		product = testhelpers.NewProductBuilder().Build()
		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(version)
	)

	// override default workflow config to empty map
	versionConfig.ObjectStores.Workflows = map[string]entity.WorkflowObjectStoresConfig{}

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrWorkflowObjectStoreNotFound)
}

func (s *StartVersionTestSuite) TestStartVersion_ErrorMapping_ProcessStreamNotFound() {
	ctx := context.Background()

	var (
		product = testhelpers.NewProductBuilder().Build()

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = s.getConfigForVersion(version)
	)

	// override default workflow config to empty map
	versionConfig.Streams.Workflows[workflow.Name] = entity.WorkflowStreamResources{
		Stream:    "stream",
		Processes: map[string]entity.ProcessStreamConfig{},
	}

	err := s.k8sVersionClient.Start(ctx, product, version, versionConfig)
	s.Assert().ErrorIs(err, entity.ErrProcessStreamNotFound)
}

func (s *StartVersionTestSuite) getConfigForVersion(version *entity.Version) *entity.VersionStreamingResources {
	var (
		workflow = version.Workflows[0]
		process  = version.Workflows[0].Processes[0]
	)

	streamConfig := &entity.VersionStreams{
		Workflows: map[string]entity.WorkflowStreamResources{
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

	objectStoreConfig := &entity.VersionObjectStores{
		Workflows: map[string]entity.WorkflowObjectStoresConfig{
			workflow.Name: {
				Processes: map[string]string{
					process.Name: "test-object-store",
				},
			},
		},
	}

	keyValueStoresConfig := &entity.KeyValueStores{
		VersionKeyValueStore: "test-product-kv-store",
		Workflows: map[string]*entity.WorkflowKeyValueStores{
			workflow.Name: {
				KeyValueStore: "test-workflow-kv-store",
				Processes: map[string]string{
					process.Name: "test-process-kv-store",
				},
			},
		},
	}

	return &entity.VersionStreamingResources{
		Streams:        streamConfig,
		ObjectStores:   objectStoreConfig,
		KeyValueStores: keyValueStoresConfig,
	}
}
