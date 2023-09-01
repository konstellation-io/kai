//go:build unit

package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
	"github.com/konstellation-io/kai/engine/k8s-manager/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterProcess(t *testing.T) {
	var (
		logger          = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		imageBuilder    = mocks.NewImageBuilderMock(t)
		processRegister = usecase.NewProcessRegister(logger, imageBuilder)
		ctx             = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
		process = "test-process"

		expectedImageRef = "test-product_test-process:v1.0.0"
	)

	imageBuilder.EXPECT().
		BuildImage(ctx, expectedImageRef, mock.Anything).
		Return(expectedImageRef, nil).
		Once()

	imageRef, err := processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		Product: product,
		Version: version,
		Process: process,
		Sources: nil,
	})
	assert.NoError(t, err)

	assert.Equal(t, expectedImageRef, imageRef)
}

func TestRegisterProcess_ErrorBuildingImage(t *testing.T) {
	var (
		logger          = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		imageBuilder    = mocks.NewImageBuilderMock(t)
		processRegister = usecase.NewProcessRegister(logger, imageBuilder)
		ctx             = context.Background()
	)

	const (
		product = "test-product"
		version = "v1.0.0"
		process = "test-process"

		expectedImageRef = "test-product_test-process:v1.0.0"
	)

	expectedError := errors.New("error building image")

	imageBuilder.EXPECT().
		BuildImage(ctx, expectedImageRef, mock.Anything).
		Return("", expectedError).
		Once()

	_, err := processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		Product: product,
		Version: version,
		Process: process,
		Sources: nil,
	})
	assert.ErrorIs(t, err, expectedError)
}
