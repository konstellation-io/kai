package k8smanager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/k8smanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/logging"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
)

var (
	productID = "test-product"
	version   = testhelpers.NewVersionBuilder().Build()
)

type StopVersionTestSuite struct {
	suite.Suite
	cfg              *config.Config
	logger           logging.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *k8smanager.K8sVersionClient
}

func TestStopVersionTestSuite(t *testing.T) {
	suite.Run(t, new(StopVersionTestSuite))
}

func (s *StopVersionTestSuite) SetupSuite() {
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

func (s *StopVersionTestSuite) TestStopVersion() {
	ctx := context.Background()

	req := &versionpb.StopRequest{
		Product: productID,
		Version: version.Name,
	}

	s.mockService.EXPECT().Stop(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Stop(ctx, productID, &version)
	s.Require().NoError(err)
}

func (s *StopVersionTestSuite) TestStopVersion_ClientError() {
	ctx := context.Background()

	expectedError := errors.New("client error")

	req := &versionpb.StopRequest{
		Product: productID,
		Version: version.Name,
	}

	s.mockService.EXPECT().Stop(gomock.Any(), req).Return(nil, expectedError)

	err := s.k8sVersionClient.Stop(ctx, productID, &version)
	s.Assert().ErrorIs(err, expectedError)
}
