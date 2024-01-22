//go:build unit

package versionservice_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/stretchr/testify/suite"
)

type VersionServiceTestSuite struct {
	suite.Suite
	logger           logr.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *versionservice.K8sVersionService
}

func TestK8sManagerTestSuite(t *testing.T) {
	suite.Run(t, new(VersionServiceTestSuite))
}

func (s *VersionServiceTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	service := mocks.NewMockVersionServiceClient(mockController)

	k8sVersionClient, err := versionservice.New(logger, service)
	s.Require().NoError(err)

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

	s.mockService.EXPECT().Publish(gomock.Any(), req).Return(&versionpb.PublishResponse{}, nil)

	err := s.k8sVersionClient.Publish(ctx, productID, version.Tag)
	s.Require().NoError(err)
}

func (s *VersionServiceTestSuite) TestPublish_ClientError() {
	ctx := context.Background()

	req := &versionpb.PublishRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	expectedError := errors.New("k8s error")

	s.mockService.EXPECT().Publish(gomock.Any(), req).Return(nil, expectedError)

	err := s.k8sVersionClient.Publish(ctx, productID, version.Tag)
	s.Require().ErrorIs(err, expectedError)
}

func (s *VersionServiceTestSuite) TestUnpublish() {
	ctx := context.Background()

	req := &versionpb.UnpublishRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	s.mockService.EXPECT().Unpublish(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Unpublish(ctx, productID, version)
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
		expectedProductID    = "test-product-id"
		expectedProcessID    = "test-process-id"
		expectedProcessImage = "test-process-image"
	)

	s.mockService.EXPECT().RegisterProcess(ctx, &versionpb.RegisterProcessRequest{
		ProductId:    expectedProductID,
		ProcessId:    expectedProcessID,
		ProcessImage: expectedProcessImage,
	}).Return(&versionpb.RegisterProcessResponse{
		ImageId: expectedProcessID,
	}, nil)

	ref, err := s.k8sVersionClient.RegisterProcess(ctx, expectedProductID, expectedProcessID, expectedProcessImage)
	s.NoError(err)
	s.Equal(expectedProcessID, ref)
}

func (s *VersionServiceTestSuite) TestRegisterProcess_ClientError() {
	ctx := context.Background()

	const (
		expectedProductID    = "test-product-id"
		expectedProcessID    = "test-process-id"
		expectedProcessImage = "test-process-image"
	)

	var (
		expectedError = errors.New("client error")
	)

	s.mockService.EXPECT().RegisterProcess(ctx, &versionpb.RegisterProcessRequest{
		ProductId:    expectedProductID,
		ProcessId:    expectedProcessID,
		ProcessImage: expectedProcessImage,
	}).Return(nil, expectedError)

	_, err := s.k8sVersionClient.RegisterProcess(ctx, expectedProductID, expectedProcessID, expectedProcessImage)
	s.ErrorIs(err, expectedError)
}
