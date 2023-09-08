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
)

func TestRegisterProcess(t *testing.T) {
	var (
		logger          = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		imageBuilder    = mocks.NewImageBuilderMock(t)
		processRegister = usecase.NewProcessRegister(logger, imageBuilder)
		ctx             = context.Background()
		source          = []byte("test-source")
	)

	const (
		product      = "test-product"
		processID    = "test-product_test-process:v1.0.0"
		processImage = "process-image"
	)

	imageBuilder.EXPECT().
		BuildImage(ctx, processID, processImage, source).
		Return(processImage, nil).
		Once()

	imageID, err := processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		ProcessID:    processID,
		ProcessImage: processImage,
		Sources:      source,
	})
	assert.NoError(t, err)

	assert.Equal(t, processImage, imageID)
}

func TestRegisterProcess_ErrorBuildingImage(t *testing.T) {
	var (
		logger          = testr.NewWithOptions(t, testr.Options{Verbosity: -1})
		imageBuilder    = mocks.NewImageBuilderMock(t)
		processRegister = usecase.NewProcessRegister(logger, imageBuilder)
		ctx             = context.Background()
		source          = []byte("test-source")
	)

	const (
		product      = "test-product"
		processID    = "test-product_test-process:v1.0.0"
		processImage = "process-image"
	)

	expectedError := errors.New("error building image")

	imageBuilder.EXPECT().
		BuildImage(ctx, processID, processImage, source).
		Return("", expectedError).
		Once()

	_, err := processRegister.RegisterProcess(ctx, usecase.RegisterProcessParams{
		ProcessID:    processID,
		ProcessImage: processImage,
		Sources:      source,
	})
	assert.ErrorIs(t, err, expectedError)
}
