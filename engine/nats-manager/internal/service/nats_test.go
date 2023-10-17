//go:build unit

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/config"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/logging"
	"github.com/konstellation-io/kai/engine/nats-manager/internal/service"
	"github.com/konstellation-io/kai/engine/nats-manager/mocks"
	"github.com/konstellation-io/kai/engine/nats-manager/proto/natspb"
	"github.com/stretchr/testify/suite"
)

const (
	productID  = "productID"
	versionTag = "v1.0.0"
)

var (
	protoWorkflows = []*natspb.Workflow{
		{
			Name: "test-workflow",
			Processes: []*natspb.Process{
				{
					Name:          "test-process",
					Subscriptions: []string{"test-subject"},
					ObjectStore: &natspb.ObjectStore{
						Name:  "test-objectStore",
						Scope: natspb.ObjectStoreScope_SCOPE_WORKFLOW,
					},
				},
			},
		},
	}

	entityWorkflows = []entity.Workflow{
		{
			Name: protoWorkflows[0].Name,
			Processes: []entity.Process{
				{
					Name:          protoWorkflows[0].Processes[0].Name,
					Subscriptions: protoWorkflows[0].Processes[0].Subscriptions,
					ObjectStore: &entity.ObjectStore{
						Name:  protoWorkflows[0].Processes[0].ObjectStore.Name,
						Scope: entity.ObjStoreScopeWorkflow,
					},
				},
			},
		},
	}
)

type NatsServiceTestSuite struct {
	suite.Suite
	cfg             *config.Config
	logger          logging.Logger
	natsManagerMock *mocks.MockNatsManager
	natsService     *service.NatsService
}

func TestNatsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NatsServiceTestSuite))
}

func (s *NatsServiceTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	cfg := &config.Config{}
	logger := mocks.NewMockLogger(mockController)
	mocks.AddLoggerExpects(logger)

	s.natsManagerMock = mocks.NewMockNatsManager(mockController)

	s.natsService = service.NewNatsService(cfg, logger, s.natsManagerMock)

	s.cfg = cfg
	s.logger = logger
}

func (s *NatsServiceTestSuite) TestCreateStreams() {
	req := &natspb.CreateStreamsRequest{
		ProductId:  productID,
		VersionTag: versionTag,
		Workflows:  protoWorkflows,
	}

	expectedEntityWorkflows := entityWorkflows

	workflowStream := "test-workflow-stream"
	processSubject := "test-process-subject"

	managerResponse := entity.WorkflowsStreamsConfig{
		req.Workflows[0].Name: &entity.StreamConfig{
			Stream: workflowStream,
			Processes: entity.ProcessesStreamConfig{
				req.Workflows[0].Processes[0].Name: entity.ProcessStreamConfig{
					Subject:       processSubject,
					Subscriptions: req.Workflows[0].Processes[0].Subscriptions,
				},
			},
		},
	}

	expectedClientResponse := &natspb.CreateStreamsResponse{
		Workflows: map[string]*natspb.WorkflowStreamConfig{
			req.Workflows[0].Name: {
				Stream: workflowStream,
				Processes: map[string]*natspb.ProcessStreamConfig{
					req.Workflows[0].Processes[0].Name: {
						Subject:       processSubject,
						Subscriptions: req.Workflows[0].Processes[0].Subscriptions,
					},
				},
			},
		},
	}

	s.natsManagerMock.EXPECT().
		CreateStreams(req.ProductId, req.VersionTag, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateStreams(nil, req)
	s.Require().NoError(err)
	s.Equal(expectedClientResponse, clientResponse)
}

func (s *NatsServiceTestSuite) TestCreateObjectStores() {
	req := &natspb.CreateObjectStoresRequest{
		ProductId:  productID,
		VersionTag: versionTag,
		Workflows:  protoWorkflows,
	}

	expectedEntityWorkflows := entityWorkflows

	testObjectStore := "test-objectStore"
	managerResponse := entity.WorkflowsObjectStoresConfig{
		req.Workflows[0].Name: &entity.WorkflowObjectStoresConfig{
			Processes: entity.ProcessesObjectStoresConfig{
				req.Workflows[0].Processes[0].Name: testObjectStore,
			},
		},
	}

	expectedClientResponse := &natspb.CreateObjectStoresResponse{
		Workflows: map[string]*natspb.WorkflowObjectStoreConfig{
			req.Workflows[0].Name: {
				Processes: map[string]string{
					req.Workflows[0].Processes[0].Name: testObjectStore,
				},
			},
		},
	}

	s.natsManagerMock.EXPECT().
		CreateObjectStores(req.ProductId, req.VersionTag, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateObjectStores(nil, req)
	s.Require().NoError(err)
	s.Equal(expectedClientResponse, clientResponse)
}

func (s *NatsServiceTestSuite) TestDeleteStreams() {
	req := &natspb.DeleteStreamsRequest{
		ProductId:  productID,
		VersionTag: versionTag,
	}

	s.natsManagerMock.EXPECT().
		DeleteStreams(req.ProductId, req.VersionTag).
		Return(nil)

	clientResponse, err := s.natsService.DeleteStreams(nil, req)
	s.Require().NoError(err)
	s.NotEmpty(clientResponse.Message)
}

func (s *NatsServiceTestSuite) TestDeleteObjectStores() {
	req := &natspb.DeleteObjectStoresRequest{
		ProductId:  productID,
		VersionTag: versionTag,
	}

	s.natsManagerMock.EXPECT().
		DeleteObjectStores(req.ProductId, req.VersionTag).
		Return(nil)

	clientResponse, err := s.natsService.DeleteObjectStores(nil, req)
	s.Require().NoError(err)
	s.NotEmpty(clientResponse.Message)
}

func (s *NatsServiceTestSuite) TestCreateVersionKeyValueStores() {
	req := &natspb.CreateVersionKeyValueStoresRequest{
		ProductId:  productID,
		VersionTag: versionTag,
		Workflows:  protoWorkflows,
	}

	expectedEntityWorkflows := entityWorkflows

	testProjectStore := "test-project-store"
	testWorkflowStore := "test-workflow-store"
	testProcessStore := "test-process-store"

	managerResponse := &entity.VersionKeyValueStores{
		ProjectStore: testProjectStore,
		WorkflowsStores: map[string]*entity.WorkflowKeyValueStores{
			req.Workflows[0].Name: {
				WorkflowStore: testWorkflowStore,
				Processes: map[string]string{
					req.Workflows[0].Processes[0].Name: testProcessStore,
				},
			},
		},
	}

	expectedClientResponse := &natspb.CreateVersionKeyValueStoresResponse{
		KeyValueStore: testProjectStore,
		Workflows: map[string]*natspb.WorkflowKeyValueStoreConfig{
			req.Workflows[0].Name: {
				KeyValueStore: testWorkflowStore,
				Processes: map[string]string{
					req.Workflows[0].Processes[0].Name: testProcessStore,
				},
			},
		},
	}

	s.natsManagerMock.EXPECT().
		CreateVersionKeyValueStores(req.ProductId, req.VersionTag, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateVersionKeyValueStores(nil, req)
	s.Require().NoError(err)
	s.Equal(expectedClientResponse, clientResponse)
}

func (s *NatsServiceTestSuite) TestCreateStreamsError() {
	req := &natspb.CreateStreamsRequest{}

	s.natsManagerMock.EXPECT().
		CreateStreams(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("mock error"))

	_, err := s.natsService.CreateStreams(nil, req)
	s.Require().Error(err)
}

func (s *NatsServiceTestSuite) TestCreateObjectStoresError() {
	req := &natspb.CreateObjectStoresRequest{}

	s.natsManagerMock.EXPECT().
		CreateObjectStores(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("mock error"))

	_, err := s.natsService.CreateObjectStores(nil, req)
	s.Require().Error(err)
}

func (s *NatsServiceTestSuite) TestDeleteStreamsError() {
	req := &natspb.DeleteStreamsRequest{}

	s.natsManagerMock.EXPECT().
		DeleteStreams(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	_, err := s.natsService.DeleteStreams(nil, req)
	s.Require().Error(err)
}

func (s *NatsServiceTestSuite) TestDeleteObjectStoresError() {
	req := &natspb.DeleteObjectStoresRequest{}

	s.natsManagerMock.EXPECT().
		DeleteObjectStores(gomock.Any(), gomock.Any()).
		Return(errors.New("mock error"))

	_, err := s.natsService.DeleteObjectStores(nil, req)
	s.Require().Error(err)
}

func (s *NatsServiceTestSuite) TestCreateVersionKeyValueStoresError() {
	req := &natspb.CreateVersionKeyValueStoresRequest{}

	s.natsManagerMock.EXPECT().
		CreateVersionKeyValueStores(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("mock error"))

	_, err := s.natsService.CreateVersionKeyValueStores(nil, req)
	s.Require().Error(err)
}

func (s *NatsServiceTestSuite) TestCreateGlobalKeyValueStores() {
	ctx := context.Background()
	req := &natspb.CreateGlobalKeyValueStoreRequest{
		ProductId: productID,
	}
	expectedKVStore := "globalKVStore"

	s.natsManagerMock.EXPECT().
		CreateGlobalKeyValueStore(productID).
		Return(expectedKVStore, nil)

	res, err := s.natsService.CreateGlobalKeyValueStore(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(res.GlobalKeyValueStore, expectedKVStore)
}

func (s *NatsServiceTestSuite) TestCreateGlobalKeyValueStores_Error() {
	ctx := context.Background()
	req := &natspb.CreateGlobalKeyValueStoreRequest{
		ProductId: productID,
	}
	expectedError := errors.New("nats manager error")

	s.natsManagerMock.EXPECT().
		CreateGlobalKeyValueStore(productID).
		Return("", expectedError)

	_, err := s.natsService.CreateGlobalKeyValueStore(ctx, req)
	s.Require().ErrorIs(err, expectedError)
}

func (s *NatsServiceTestSuite) TestUpdateKeyValueConfiguration() {
	ctx := context.Background()
	kvStore := "kvStore"
	configuration := map[string]string{
		"key1": "val1",
	}
	req := &natspb.UpdateKeyValueConfigurationRequest{
		KeyValueStoresConfig: []*natspb.KeyValueConfiguration{
			{
				KeyValueStore: kvStore,
				Configuration: configuration,
			},
		},
	}

	s.natsManagerMock.EXPECT().UpdateKeyValueStoresConfiguration([]entity.KeyValueConfiguration{
		{
			KeyValueStore: kvStore,
			Configuration: configuration,
		},
	}).Return(nil)

	_, err := s.natsService.UpdateKeyValueConfiguration(ctx, req)
	s.Require().NoError(err)
}

func (s *NatsServiceTestSuite) TestUpdateKeyValueConfiguration_Error() {
	ctx := context.Background()
	req := &natspb.UpdateKeyValueConfigurationRequest{
		KeyValueStoresConfig: []*natspb.KeyValueConfiguration{
			{
				KeyValueStore: "kvStore",
			},
		},
	}
	expectedError := errors.New("nats manager error")

	s.natsManagerMock.EXPECT().UpdateKeyValueStoresConfiguration(gomock.Any()).Return(expectedError)

	_, err := s.natsService.UpdateKeyValueConfiguration(ctx, req)
	s.Require().ErrorIs(err, expectedError)
}
