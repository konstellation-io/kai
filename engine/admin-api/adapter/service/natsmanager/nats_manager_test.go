//go:build unit

package natsmanager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/natsmanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

const productID = "test-product"

var (
	testProcess = testhelpers.NewProcessBuilder().
			WithObjectStore(&entity.ProcessObjectStore{
			Name:  "test-object-store",
			Scope: entity.ObjectStoreScopeWorkflow,
		}).
		WithNetworking(&entity.ProcessNetworking{
			TargetPort:      8080,
			DestinationPort: 8080,
			Protocol:        "GRPC",
		}).
		Build()

	testWorkflow = testhelpers.NewWorkflowBuilder().
			WithProcesses([]entity.Process{testProcess}).
			Build()

	testVersion = testhelpers.NewVersionBuilder().
			WithWorkflows([]entity.Workflow{testWorkflow}).
			Build()

	testReqWorkflows = []*natspb.Workflow{
		{
			Name: testWorkflow.Name,
			Processes: []*natspb.Process{
				{
					Name: testProcess.Name,
					ObjectStore: &natspb.ObjectStore{
						Name:  testProcess.ObjectStore.Name,
						Scope: natspb.ObjectStoreScope_SCOPE_WORKFLOW,
					},
					Subscriptions: testProcess.Subscriptions,
				},
			},
		},
	}
)

type NatsManagerTestSuite struct {
	suite.Suite
	cfg               *config.Config
	logger            logging.Logger
	mockService       *mocks.MockNatsManagerServiceClient
	natsManagerClient *natsmanager.Client
}

func TestNatsManagerTestSuite(t *testing.T) {
	suite.Run(t, new(NatsManagerTestSuite))
}

func (s *NatsManagerTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	cfg := &config.Config{}
	logger := mocks.NewMockLogger(mockController)
	mocks.AddLoggerExpects(logger)
	service := mocks.NewMockNatsManagerServiceClient(mockController)

	k8sVersionClient, err := natsmanager.NewClient(cfg, logger, service)
	s.Require().NoError(err)

	s.cfg = cfg
	s.logger = logger
	s.mockService = service
	s.natsManagerClient = k8sVersionClient
}

func (s *NatsManagerTestSuite) TestCreateStreams() {
	ctx := context.Background()

	req := &natspb.CreateStreamsRequest{
		ProductId:  productID,
		VersionTag: testVersion.Tag,
		Workflows:  testReqWorkflows,
	}

	natsManagerResponse := &natspb.CreateStreamsResponse{
		Workflows: map[string]*natspb.WorkflowStreamConfig{
			testWorkflow.Name: {
				Stream: "test-workflow-stream-name",
				Processes: map[string]*natspb.ProcessStreamConfig{
					testProcess.Name: {
						Subject:       "test-process-subject-name",
						Subscriptions: testProcess.Subscriptions,
					},
				},
			},
		},
	}

	expctedResponse := &entity.VersionStreams{
		Workflows: map[string]entity.WorkflowStreamResources{
			testWorkflow.Name: {
				Stream: "test-workflow-stream-name",
				Processes: map[string]entity.ProcessStreamConfig{
					testProcess.Name: {
						Subject:       "test-process-subject-name",
						Subscriptions: testProcess.Subscriptions,
					},
				},
			},
		},
	}

	s.mockService.EXPECT().CreateStreams(ctx, req).Return(natsManagerResponse, nil)

	res, err := s.natsManagerClient.CreateStreams(ctx, productID, testVersion)
	s.Require().NoError(err)
	s.Equal(expctedResponse, res)
}

func (s *NatsManagerTestSuite) TestCreateObjectStores() {
	ctx := context.Background()

	req := &natspb.CreateObjectStoresRequest{
		ProductId:  productID,
		VersionTag: testVersion.Tag,
		Workflows:  testReqWorkflows,
	}

	natsManagerResponse := &natspb.CreateObjectStoresResponse{
		Workflows: map[string]*natspb.WorkflowObjectStoreConfig{
			testWorkflow.Name: {
				Processes: map[string]string{
					testProcess.Name: "test-object-store-name",
				},
			},
		},
	}

	expectedResponse := &entity.VersionObjectStores{
		Workflows: map[string]entity.WorkflowObjectStoresConfig{
			testWorkflow.Name: {
				Processes: map[string]string{
					testProcess.Name: "test-object-store-name",
				},
			},
		},
	}

	s.mockService.EXPECT().CreateObjectStores(ctx, req).Return(natsManagerResponse, nil)

	res, err := s.natsManagerClient.CreateObjectStores(ctx, productID, testVersion)
	s.Require().NoError(err)
	s.Equal(expectedResponse, res)
}

func (s *NatsManagerTestSuite) TestCreateKeyValueStores() {
	ctx := context.Background()

	req := &natspb.CreateVersionKeyValueStoresRequest{
		ProductId:  productID,
		VersionTag: testVersion.Tag,
		Workflows:  testReqWorkflows,
	}

	natsManagerResponse := &natspb.CreateVersionKeyValueStoresResponse{
		KeyValueStore: "v1.0.0-key-value-store-name",
		Workflows: map[string]*natspb.WorkflowKeyValueStoreConfig{
			testWorkflow.Name: {
				KeyValueStore: "test-workflow-key-value-store-name",
				Processes: map[string]string{
					testProcess.Name: "test-process-key-value-store-name",
				},
			},
		},
	}

	expectedResponse := &entity.KeyValueStores{
		VersionKeyValueStore: "v1.0.0-key-value-store-name",
		Workflows: map[string]*entity.WorkflowKeyValueStores{
			testWorkflow.Name: {
				KeyValueStore: "test-workflow-key-value-store-name",
				Processes: map[string]string{
					testProcess.Name: "test-process-key-value-store-name",
				},
			},
		},
	}

	s.mockService.EXPECT().CreateVersionKeyValueStores(ctx, req).Return(natsManagerResponse, nil)

	res, err := s.natsManagerClient.CreateVersionKeyValueStores(ctx, productID, testVersion)
	s.Require().NoError(err)
	s.Equal(expectedResponse, res)
}

func (s *NatsManagerTestSuite) TestDeleteStreams() {
	ctx := context.Background()

	req := &natspb.DeleteStreamsRequest{
		ProductId:  productID,
		VersionTag: testVersion.Tag,
	}

	s.mockService.EXPECT().DeleteStreams(ctx, req).Return(&natspb.DeleteResponse{}, nil)

	err := s.natsManagerClient.DeleteStreams(ctx, productID, testVersion.Tag)
	s.Require().NoError(err)
}

func (s *NatsManagerTestSuite) TestDeleteObjectStores() {
	ctx := context.Background()

	req := &natspb.DeleteObjectStoresRequest{
		ProductId:  productID,
		VersionTag: testVersion.Tag,
	}

	s.mockService.EXPECT().DeleteObjectStores(ctx, req).Return(&natspb.DeleteResponse{}, nil)

	err := s.natsManagerClient.DeleteObjectStores(ctx, productID, testVersion.Tag)
	s.Require().NoError(err)
}

func (s *NatsManagerTestSuite) TestCreateStreamsManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().CreateStreams(ctx, gomock.Any()).
		Return(&natspb.CreateStreamsResponse{}, errors.New("mocked error"))

	_, err := s.natsManagerClient.CreateStreams(ctx, productID, testVersion)
	s.Error(err)
}

func (s *NatsManagerTestSuite) TestCreateObjectStoresManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().CreateObjectStores(ctx, gomock.Any()).
		Return(&natspb.CreateObjectStoresResponse{}, errors.New("mocked error"))

	_, err := s.natsManagerClient.CreateObjectStores(ctx, productID, testVersion)
	s.Error(err)
}

func (s *NatsManagerTestSuite) TestCreateKeyValueStoresManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().CreateVersionKeyValueStores(ctx, gomock.Any()).
		Return(&natspb.CreateVersionKeyValueStoresResponse{}, errors.New("mocked error"))

	_, err := s.natsManagerClient.CreateVersionKeyValueStores(ctx, productID, testVersion)
	s.Error(err)
}

func (s *NatsManagerTestSuite) TestDeleteStreamsManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().DeleteStreams(ctx, gomock.Any()).
		Return(&natspb.DeleteResponse{}, errors.New("mocked error"))

	err := s.natsManagerClient.DeleteStreams(ctx, productID, testVersion.Tag)
	s.Error(err)
}

func (s *NatsManagerTestSuite) TestDeleteObjectStoresManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().DeleteObjectStores(ctx, gomock.Any()).
		Return(&natspb.DeleteResponse{}, errors.New("mocked error"))

	err := s.natsManagerClient.DeleteObjectStores(ctx, productID, testVersion.Tag)
	s.Error(err)
}
