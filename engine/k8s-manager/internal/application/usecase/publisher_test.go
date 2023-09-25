//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishVersion(t *testing.T) {
	var (
		ctx              = context.Background()
		logger           = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		networkPublisher = mocks.NewContainerServiceMock(t)
		versionPublisher = usecase.NewVersionPublisher(logger, networkPublisher)

		product = "test-product"
		version = "v1.0.0"
	)

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	networkPublisher.EXPECT().PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	}).Return(expectedURLs, nil)

	urls, err := versionPublisher.PublishVersion(ctx, product, version)
	require.NoError(t, err)
	assert.Equal(t, expectedURLs, urls)
}

func TestPublishVersion_Error(t *testing.T) {
	var (
		ctx              = context.Background()
		logger           = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		networkPublisher = mocks.NewContainerServiceMock(t)
		versionPublisher = usecase.NewVersionPublisher(logger, networkPublisher)

		product = "test-product"
		version = "v1.0.0"
	)

	expectedError := errors.New("publish network urls")

	networkPublisher.EXPECT().PublishNetwork(ctx, service.PublishNetworkParams{
		Product: product,
		Version: version,
	}).Return(nil, expectedError)

	_, err := versionPublisher.PublishVersion(ctx, product, version)
	require.Error(t, expectedError, err)
}
