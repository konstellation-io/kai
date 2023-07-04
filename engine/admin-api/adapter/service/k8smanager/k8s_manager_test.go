//go:build unit

package k8smanager_test

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/k8smanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/suite"
)

type K8sManagerTestSuite struct {
	suite.Suite
	cfg              *config.Config
	logger           logging.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *k8smanager.K8sVersionClient
}

func TestK8sManagerTestSuite(t *testing.T) {
	suite.Run(t, new(K8sManagerTestSuite))
}

func (s *K8sManagerTestSuite) SetupSuite() {
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

func (s *K8sManagerTestSuite) TestPublish() {
	ctx := context.Background()

	req := &versionpb.PublishRequest{
		Product: productID,
		Version: version.Name,
	}

	s.mockService.EXPECT().Publish(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Publish(ctx, productID, &version)
	s.Require().NoError(err)
}

func (s *K8sManagerTestSuite) TestUnpublish() {
	ctx := context.Background()

	req := &versionpb.UnpublishRequest{
		Product: productID,
		Version: version.Name,
	}

	s.mockService.EXPECT().Unpublish(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Unpublish(ctx, productID, &version)
	s.Require().NoError(err)
}

func (s *K8sManagerTestSuite) TestWatchProcessStatus() {
	ctx := context.Background()

	req := &versionpb.ProcessStatusRequest{
		ProductId:   productID,
		VersionName: version.Name,
	}

	processStatusResponse := &versionpb.ProcessStatusResponse{
		ProcessId: "test-process-id",
		Name:      "test-process-name",
		Status:    "STARTED",
	}

	stream := mocks.NewMockVersionService_WatchProcessStatusClient(gomock.NewController(s.T()))

	s.mockService.EXPECT().WatchProcessStatus(ctx, req).Return(stream, nil)
	stream.EXPECT().Recv().Return(processStatusResponse, nil)
	stream.EXPECT().Context().Return(ctx).Times(2)
	stream.EXPECT().Recv().Return(nil, io.EOF)

	statusChannel, err := s.k8sVersionClient.WatchProcessStatus(ctx, productID, version.Name)
	s.Require().NoError(err)

	process := <-statusChannel
	s.Require().Equal(entity.ProcessStatusStarted, process.Status)

	process = <-statusChannel
	s.Nil(process)
}
