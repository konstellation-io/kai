package k8smanager_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/config"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/k8smanager"
	"github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	"github.com/konstellation-io/kai/engine/admin-api/mocks"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestStopVersion(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		cfg     = &config.Config{}
		logger  = mocks.NewMockLogger(ctrl)
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID = "test-product"
		version   = testhelpers.NewVersionBuilder().Build()
	)

	req := &versionpb.StopRequest{
		Product: productID,
		Version: version.Name,
	}

	service.EXPECT().Stop(gomock.Any(), req).Return(&versionpb.Response{Message: "ok"}, nil)

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Stop(ctx, productID, &version)
	assert.NoError(t, err)
}

func TestStopVersion_ClientError(t *testing.T) {
	var (
		ctrl    = gomock.NewController(t)
		service = mocks.NewMockVersionServiceClient(ctrl)
		cfg     = &config.Config{}
		logger  = mocks.NewMockLogger(ctrl)
		ctx     = context.Background()
	)

	mocks.AddLoggerExpects(logger)

	var (
		productID = "test-product"
		version   = testhelpers.NewVersionBuilder().Build()
	)

	expectedError := errors.New("client error")

	req := &versionpb.StopRequest{
		Product: productID,
		Version: version.Name,
	}

	service.EXPECT().Stop(gomock.Any(), req).Return(nil, expectedError)

	client, _ := k8smanager.NewK8sVersionClient(cfg, logger, service)

	err := client.Stop(ctx, productID, &version)
	assert.ErrorIs(t, err, expectedError)
}
