//go:build unit

package process_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
)

func (s *ProcessHandlerTestSuite) TestDeleteProcess_WithProduct() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteRegisteredProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(nil)
	s.processRepo.EXPECT().Delete(ctx, opts.Product, processID).Return(nil)

	returnedProcessID, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.Require().NoError(err)

	s.Equal(processID, returnedProcessID)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_Public() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: true,
	}

	imageName := "kai_process-name"
	processID := "kai_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActDeletePublicProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, _publicRegistry, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(nil)
	s.processRepo.EXPECT().Delete(ctx, _publicRegistry, processID).Return(nil)

	returnedProcessID, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.Require().NoError(err)

	s.Equal(processID, returnedProcessID)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_MissingProductInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{}

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, process.ErrMissingProductInParams)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_MissingVersionInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product: "test-product",
	}

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, process.ErrMissingVersionInParams)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_MissingProcessInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product: "test-product",
		Version: "v1.0.0",
	}

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, process.ErrMissingProcessInParams)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_IsPublicAndHasProduct() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: true,
	}

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, process.ErrIsPublicAndHasProduct)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_NoProductGrants() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}
	expectedErr := errors.New("auth error")

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteRegisteredProcess).Return(expectedErr)

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, expectedErr)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_NoRoleGrants() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: true,
	}
	expectedErr := errors.New("auth error")

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActDeletePublicProcess).Return(expectedErr)

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, expectedErr)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_GetByIDError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteRegisteredProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, process.ErrRegisteredProcessNotFound)

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, process.ErrRegisteredProcessNotFound)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_DeleteProcessError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}
	expectedErr := errors.New("delete process error")

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteRegisteredProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(expectedErr)

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, expectedErr)
}

func (s *ProcessHandlerTestSuite) TestDeleteProcess_DeleteProcessRepoError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}
	expectedErr := errors.New("delete process error")

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteRegisteredProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(nil)
	s.processRepo.EXPECT().Delete(ctx, opts.Product, processID).Return(expectedErr)

	_, err := s.processHandler.DeleteProcess(ctx, user, opts)
	s.ErrorIs(err, expectedErr)
}
