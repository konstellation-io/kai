//go:build unit

package process_test

import (
	"context"
	"fmt"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
)

func (s *ProcessServiceTestSuite) TestDeleteProcess_WithProduct() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(nil)
	s.processRepo.EXPECT().Delete(ctx, opts.Product, processID).Return(nil)

	returnedProcessID, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().NoError(err)

	s.Equal(processID, returnedProcessID)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_Public() {
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

	returnedProcessID, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().NoError(err)

	s.Equal(processID, returnedProcessID)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_MissingProductInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{}

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(process.ErrMissingProductInParams, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_MissingVersionInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product: "test-product",
	}

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(process.ErrMissingVersionInParams, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_MissingProcessInDeleteOptions() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product: "test-product",
		Version: "v1.0.0",
	}

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(process.ErrMissingProcessInParams, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_IsPublicAndHasProduct() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: true,
	}

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(process.ErrIsPublicAndHasProduct, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_NoProductGrants() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(auth.UnauthorizedError{})

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(auth.UnauthorizedError{}, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_NoRoleGrants() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: true,
	}

	s.accessControl.EXPECT().CheckRoleGrants(user, auth.ActDeletePublicProcess).Return(auth.UnauthorizedError{})

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(auth.UnauthorizedError{}, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_GetByIDError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, process.ErrRegisteredProcessNotFound)

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
	s.Equal(process.ErrRegisteredProcessNotFound, err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_DeleteProcessError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(fmt.Errorf("error"))

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
}

func (s *ProcessServiceTestSuite) TestDeleteProcess_DeleteProcessRepoError() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	imageName := "test-product_process-name"
	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRegistry.EXPECT().DeleteProcess(ctx, imageName, opts.Version).Return(nil)
	s.processRepo.EXPECT().Delete(ctx, opts.Product, processID).Return(fmt.Errorf("error"))

	_, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().Error(err)
}
