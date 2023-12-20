//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/mocks"
	"github.com/stretchr/testify/require"
)

func TestUnpublishVersion(t *testing.T) {
	var (
		ctx                = context.Background()
		logger             = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		networkUnpublisher = mocks.NewContainerServiceMock(t)
		versionUnpublisher = usecase.NewVersionUnpublisher(logger, networkUnpublisher)

		product = "test-product"
		version = "v1.0.0"
	)

	networkUnpublisher.EXPECT().UnpublishNetwork(ctx,
		product,
		version,
	).Return(nil)

	err := versionUnpublisher.UnpublishVersion(ctx, product, version)
	require.NoError(t, err)
}

func TestUnpublishVersion_Error(t *testing.T) {
	var (
		ctx                = context.Background()
		logger             = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		networkUnpublisher = mocks.NewContainerServiceMock(t)
		versionUnpublisher = usecase.NewVersionUnpublisher(logger, networkUnpublisher)

		product = "test-product"
		version = "v1.0.0"
	)

	expectedError := errors.New("unpublish network error")

	networkUnpublisher.EXPECT().UnpublishNetwork(ctx,
		product,
		version,
	).Return(expectedError)

	err := versionUnpublisher.UnpublishVersion(ctx, product, version)
	require.Error(t, err)
}
