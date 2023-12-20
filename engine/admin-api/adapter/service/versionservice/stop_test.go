//go:build unit

package versionservice_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/versionservice"
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
	logger           logr.Logger
	mockService      *mocks.MockVersionServiceClient
	k8sVersionClient *versionservice.K8sVersionService
}

func TestStopVersionTestSuite(t *testing.T) {
	suite.Run(t, new(StopVersionTestSuite))
}

func (s *StopVersionTestSuite) SetupSuite() {
	mockController := gomock.NewController(s.T())
	logger := testr.NewWithOptions(s.T(), testr.Options{Verbosity: -1})
	service := mocks.NewMockVersionServiceClient(mockController)

	k8sVersionClient, err := versionservice.New(logger, service)
	s.Require().NoError(err)

	s.logger = logger
	s.mockService = service
	s.k8sVersionClient = k8sVersionClient
}

func (s *StopVersionTestSuite) TestStopVersion() {
	ctx := context.Background()

	req := &versionpb.StopRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	s.mockService.EXPECT().Stop(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	err := s.k8sVersionClient.Stop(ctx, productID, version)
	s.Require().NoError(err)
}

func (s *StopVersionTestSuite) TestStopVersion_ClientError() {
	ctx := context.Background()

	expectedError := errors.New("client error")

	req := &versionpb.StopRequest{
		Product:    productID,
		VersionTag: version.Tag,
	}

	s.mockService.EXPECT().Stop(gomock.Any(), req).Return(nil, expectedError)

	err := s.k8sVersionClient.Stop(ctx, productID, version)
	s.Assert().ErrorIs(err, expectedError)
}
