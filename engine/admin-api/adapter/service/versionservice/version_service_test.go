//go:build unit

package versionservice_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/suite"
)

type VersionServiceTestSuite struct {
	suite.Suite
	cfg              *config.Config
	logger           logging.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *versionservice.K8sVersionService
}

func TestK8sManagerTestSuite(t *testing.T) {
	suite.Run(t, new(VersionServiceTestSuite))
}

func (s *VersionServiceTestSuite) SetupSuite() {
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

func (s *VersionServiceTestSuite) TestPublish() {
	ctx := context.Background()

	req := &versionpb.PublishRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	s.mockService.EXPECT().Publish(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Publish(ctx, productID, &version)
	s.Require().NoError(err)
}

func (s *VersionServiceTestSuite) TestUnpublish() {
	ctx := context.Background()

	req := &versionpb.UnpublishRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	s.mockService.EXPECT().Unpublish(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Unpublish(ctx, productID, &version)
	s.Require().NoError(err)
}

func (s *VersionServiceTestSuite) TestWatchProcessStatus() {
	ctx := context.Background()

	req := &versionpb.ProcessStatusRequest{
		ProductId:  productID,
		VersionTag: version.Tag,
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

	statusChannel, err := s.k8sVersionClient.WatchProcessStatus(ctx, productID, version.Tag)
	s.Require().NoError(err)

	process := <-statusChannel
	s.Require().Equal(entity.ProcessStatusStarted, process.Status)

	process = <-statusChannel
	s.Nil(process)
}

func (s *VersionServiceTestSuite) TestWatchProcessStatusManagerError() {
	ctx := context.Background()

	s.mockService.EXPECT().WatchProcessStatus(ctx, gomock.Any()).Return(nil, errors.New("mocked error"))

	statusChannel, err := s.k8sVersionClient.WatchProcessStatus(ctx, productID, version.Tag)
	s.Error(err)
	s.Nil(statusChannel)
}

func (s *VersionServiceTestSuite) TestWatchProcessStatusContextCancelled() {
	ctx := context.Background()
	cancelContext, cancel := context.WithCancel(ctx)

	processStatusResponse := &versionpb.ProcessStatusResponse{}

	stream := mocks.NewMockVersionService_WatchProcessStatusClient(gomock.NewController(s.T()))

	s.mockService.EXPECT().WatchProcessStatus(cancelContext, gomock.Any()).Return(stream, nil)
	stream.EXPECT().Recv().Return(processStatusResponse, nil)
	stream.EXPECT().Context().Return(cancelContext)

	cancel()
	statusChannel, err := s.k8sVersionClient.WatchProcessStatus(cancelContext, productID, version.Tag)
	s.Require().NoError(err)

	process := <-statusChannel
	s.Nil(process)
}

func (s *VersionServiceTestSuite) TestWatchProcessStatusUnexpectedError() {
	ctx := context.Background()

	processStatusResponse := &versionpb.ProcessStatusResponse{}

	stream := mocks.NewMockVersionService_WatchProcessStatusClient(gomock.NewController(s.T()))

	s.mockService.EXPECT().WatchProcessStatus(ctx, gomock.Any()).Return(stream, nil)
	stream.EXPECT().Recv().Return(processStatusResponse, errors.New("mocked error"))
	stream.EXPECT().Context().Return(ctx)

	statusChannel, err := s.k8sVersionClient.WatchProcessStatus(ctx, productID, version.Tag)
	s.Require().NoError(err)

	process := <-statusChannel
	s.Nil(process)
}

func (s *VersionServiceTestSuite) TestRegisterProcess() {
	ctx := context.Background()

	const (
		expectedProcessID    = "test-process-id"
		expectedProcessImage = "test-process-image"
	)

	var file []byte

	s.mockService.EXPECT().RegisterProcess(ctx, &versionpb.RegisterProcessRequest{
		ProcessId:    expectedProcessID,
		ProcessImage: expectedProcessImage,
		File:         file,
	}).Return(&versionpb.RegisterProcessResponse{
		ImageId: expectedProcessID,
	}, nil)

	ref, err := s.k8sVersionClient.RegisterProcess(ctx, expectedProcessID, expectedProcessImage, file)
	s.NoError(err)
	s.Equal(expectedProcessID, ref)
}

func (s *VersionServiceTestSuite) TestRegisterProcess_ClientError() {
	ctx := context.Background()

	const (
		expectedProcessID    = "test-process-id"
		expectedProcessImage = "test-process-image"
	)

	var (
		file          []byte
		expectedError = errors.New("client error")
	)

	s.mockService.EXPECT().RegisterProcess(ctx, &versionpb.RegisterProcessRequest{
		ProcessId:    expectedProcessID,
		ProcessImage: expectedProcessImage,
		File:         file,
	}).Return(nil, expectedError)

	_, err := s.k8sVersionClient.RegisterProcess(ctx, expectedProcessID, expectedProcessImage, file)
	s.ErrorIs(err, expectedError)
}
