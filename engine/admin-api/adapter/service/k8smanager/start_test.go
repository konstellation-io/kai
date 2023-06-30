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
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
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

func TestStartVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockVersionServiceClient(ctrl)

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

		versionConfig = getConfigForVersion(&version)

		workflowStreamCfg = versionConfig.StreamsConfig.Workflows[workflow.Name]
		processStreamCfg  = workflowStreamCfg.Processes[process.Name]

		workflowKVStoreCfg = versionConfig.KeyValueStoresConfig.Workflows[workflow.Name]
		processKVStore     = workflowKVStoreCfg.Processes[process.Name]
	)

	req := &versionpb.StartRequest{
		ProductId:     productID,
		VersionName:   version.Name,
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
					},
				},
			},
		},
	}

	customMatcher := newStartRequestMatcher(req)
	service.EXPECT().Start(ctx, customMatcher).Return(&versionpb.Response{Message: "ok"}, nil)

	cfg := &config.Config{}
	logger := mocks.NewMockLogger(ctrl)
	mocks.AddLoggerExpects(logger)

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Start(ctx, productID, &version, versionConfig)

	assert.NoError(t, err)
}

func TestStartVersion_ClientError(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		logger  = mocks.NewMockLogger(ctrl)
		cfg     = &config.Config{}
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID     = "test-product"
		version       = testhelpers.NewVersionBuilder().Build()
		versionConfig = getConfigForVersion(&version)
	)

	expectedError := errors.New("client error")

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)
	service.EXPECT().Start(gomock.Any(), gomock.Any()).Return(nil, expectedError)

	err := client.Start(ctx, productID, &version, versionConfig)
	assert.ErrorIs(t, err, expectedError)
}

func TestStartVersion_ErrorMapping_WorkflowStreamFound(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		logger  = mocks.NewMockLogger(ctrl)
		cfg     = &config.Config{}
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

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

		versionConfig = getConfigForVersion(&version)
	)

	// override default workflow to empty map
	versionConfig.StreamsConfig.Workflows = map[string]entity.WorkflowStreamConfig{}

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Start(ctx, productID, &version, versionConfig)
	assert.ErrorIs(t, err, entity.ErrWorkflowStreamNotFound)
}

func TestStartVersion_ErrorMapping_NoWorkflowKeyValueStoreFound(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		logger  = mocks.NewMockLogger(ctrl)
		cfg     = &config.Config{}
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.KeyValueStoresConfig.Workflows = entity.WorkflowsKeyValueStoresConfig{}

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Start(ctx, productID, &version, versionConfig)
	assert.ErrorIs(t, err, entity.ErrWorkflowKVStoreNotFound)
}

func TestStartVersion_ErrorMapping_NoWorkflowObjectStoreFound(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		logger  = mocks.NewMockLogger(ctrl)
		cfg     = &config.Config{}
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.ObjectStoresConfig.Workflows = map[string]entity.WorkflowObjectStoresConfig{}

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Start(ctx, productID, &version, versionConfig)
	assert.ErrorIs(t, err, entity.ErrWorkflowObjectStoreNotFound)
}

func TestStartVersion_ErrorMapping_ProcessStreamNotFound(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		logger  = mocks.NewMockLogger(ctrl)
		cfg     = &config.Config{}
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID = "test-product"

		process = testhelpers.NewProcessBuilder().Build()

		workflow = testhelpers.NewWorkflowBuilder().
				WithProcesses([]entity.Process{process}).
				Build()

		version = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{workflow}).
			Build()

		versionConfig = getConfigForVersion(&version)
	)

	// override default workflow config to empty map
	versionConfig.StreamsConfig.Workflows[workflow.Name] = entity.WorkflowStreamConfig{
		Stream:    "stream",
		Processes: map[string]entity.ProcessStreamConfig{},
	}

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Start(ctx, productID, &version, versionConfig)
	assert.ErrorIs(t, err, entity.ErrProcessStreamNotFound)
}

func getConfigForVersion(version *entity.Version) *entity.VersionConfig {
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
