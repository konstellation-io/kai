//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/service/mocks"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/stretchr/testify/assert"
)

func TestStopVersion(t *testing.T) {
	var (
		logger       = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		containerSvc = mocks.NewContainerService(t)
		stopper      = usecase.NewVersionStopper(logger, containerSvc)
		ctx          = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
	)

	containerSvc.EXPECT().
		DeleteProcesses(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteConfiguration(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteNetwork(ctx, product, version).
		Return(nil).
		Once()

	err := stopper.StopVersion(ctx, usecase.StopParams{
		Product: product,
		Version: version,
	})
	assert.NoError(t, err)
}

func TestStopVersion_ErrorDeletingConfiguration(t *testing.T) {
	var (
		logger       = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		containerSvc = mocks.NewContainerService(t)
		stopper      = usecase.NewVersionStopper(logger, containerSvc)
		ctx          = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
	)

	expectedErr := errors.New("error deleting configuration")

	containerSvc.EXPECT().
		DeleteConfiguration(ctx, product, version).
		Return(expectedErr).
		Once()

	containerSvc.EXPECT().
		DeleteNetwork(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteProcesses(ctx, product, version).
		Return(nil).
		Once()

	err := stopper.StopVersion(ctx, usecase.StopParams{
		Product: product,
		Version: version,
	})
	assert.ErrorIs(t, err, expectedErr)
}

func TestStopVersion_ErrorDeletingNetwork(t *testing.T) {
	var (
		logger       = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		containerSvc = mocks.NewContainerService(t)
		stopper      = usecase.NewVersionStopper(logger, containerSvc)
		ctx          = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
	)

	expectedErr := errors.New("error deleting network")

	containerSvc.EXPECT().
		DeleteConfiguration(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteNetwork(ctx, product, version).
		Return(expectedErr).
		Once()

	containerSvc.EXPECT().
		DeleteProcesses(ctx, product, version).
		Return(nil).
		Once()

	err := stopper.StopVersion(ctx, usecase.StopParams{
		Product: product,
		Version: version,
	})
	assert.ErrorIs(t, err, expectedErr)
}

func TestStopVersion_ErrorDeletingProcesses(t *testing.T) {
	var (
		logger       = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		containerSvc = mocks.NewContainerService(t)
		stopper      = usecase.NewVersionStopper(logger, containerSvc)
		ctx          = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
	)

	expectedErr := errors.New("error deleting processes")

	containerSvc.EXPECT().
		DeleteConfiguration(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteNetwork(ctx, product, version).
		Return(nil).
		Once()

	containerSvc.EXPECT().
		DeleteProcesses(ctx, product, version).
		Return(expectedErr).
		Once()

	err := stopper.StopVersion(ctx, usecase.StopParams{
		Product: product,
		Version: version,
	})
	assert.ErrorIs(t, err, expectedErr)
}

func TestStopVersion_ErrorDeletingAllResources(t *testing.T) {
	var (
		logger       = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		containerSvc = mocks.NewContainerService(t)
		stopper      = usecase.NewVersionStopper(logger, containerSvc)
		ctx          = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
	)

	var (
		expectedDeleteConfigurationErr = errors.New("error deleting configuration")
		expectedDeleteNetworkErr       = errors.New("error deleting network")
		expectedDeleteProcessesErr     = errors.New("error deleting processes")
	)

	containerSvc.EXPECT().
		DeleteConfiguration(ctx, product, version).
		Return(expectedDeleteConfigurationErr).
		Once()

	containerSvc.EXPECT().
		DeleteNetwork(ctx, product, version).
		Return(expectedDeleteNetworkErr).
		Once()

	containerSvc.EXPECT().
		DeleteProcesses(ctx, product, version).
		Return(expectedDeleteProcessesErr).
		Once()

	err := stopper.StopVersion(ctx, usecase.StopParams{
		Product: product,
		Version: version,
	})
	assert.ErrorIs(t, err, expectedDeleteConfigurationErr)
	assert.ErrorIs(t, err, expectedDeleteNetworkErr)
	assert.ErrorIs(t, err, expectedDeleteProcessesErr)
}
