//go:build unit

package service_test

import (
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
	productID   = "productID"
	versionName = "versionName"
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
		ProductId:   productID,
		VersionName: versionName,
		Workflows:   protoWorkflows,
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
		CreateStreams(req.ProductId, req.VersionName, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateStreams(nil, req)
	s.Require().NoError(err)
	s.Equal(expectedClientResponse, clientResponse)
}

func (s *NatsServiceTestSuite) TestCreateObjectStores() {
	req := &natspb.CreateObjectStoresRequest{
		ProductId:   productID,
		VersionName: versionName,
		Workflows:   protoWorkflows,
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
		CreateObjectStores(req.ProductId, req.VersionName, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateObjectStores(nil, req)
	s.Require().NoError(err)
	s.Equal(expectedClientResponse, clientResponse)
}

func (s *NatsServiceTestSuite) TestDeleteStreams() {
	req := &natspb.DeleteStreamsRequest{
		ProductId:   productID,
		VersionName: versionName,
	}

	s.natsManagerMock.EXPECT().
		DeleteStreams(req.ProductId, req.VersionName).
		Return(nil)

	clientResponse, err := s.natsService.DeleteStreams(nil, req)
	s.Require().NoError(err)
	s.NotEmpty(clientResponse.Message)
}

func (s *NatsServiceTestSuite) TestDeleteObjectStores() {
	req := &natspb.DeleteObjectStoresRequest{
		ProductId:   productID,
		VersionName: versionName,
	}

	s.natsManagerMock.EXPECT().
		DeleteObjectStores(req.ProductId, req.VersionName).
		Return(nil)

	clientResponse, err := s.natsService.DeleteObjectStores(nil, req)
	s.Require().NoError(err)
	s.NotEmpty(clientResponse.Message)
}

func (s *NatsServiceTestSuite) TestCreateKeyValueStores() {
	req := &natspb.CreateKeyValueStoresRequest{
		ProductId:   productID,
		VersionName: versionName,
		Workflows:   protoWorkflows,
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

	expectedClientResponse := &natspb.CreateKeyValueStoreResponse{
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
		CreateKeyValueStores(req.ProductId, req.VersionName, expectedEntityWorkflows).
		Return(managerResponse, nil)

	clientResponse, err := s.natsService.CreateKeyValueStores(nil, req)
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

func (s *NatsServiceTestSuite) TestCreateKeyValueStoresError() {
	req := &natspb.CreateKeyValueStoresRequest{}

	s.natsManagerMock.EXPECT().
		CreateKeyValueStores(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, errors.New("mock error"))

	_, err := s.natsService.CreateKeyValueStores(nil, req)
	s.Require().Error(err)
}
