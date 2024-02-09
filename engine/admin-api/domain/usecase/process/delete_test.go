//go:build unit

package process_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/process"
)

// TODO
// fix this, deleteImageTag doesnt work in unitary tests, rather mock it or integration test idk
func (s *ProcessServiceTestSuite) TestDeleteWithProduct() {
	ctx := context.Background()

	opts := process.DeleteProcessOpts{
		Product:  "test-product",
		Version:  "v1.0.0",
		Process:  "process-name",
		IsPublic: false,
	}

	processID := "test-product_process-name:v1.0.0"

	s.accessControl.EXPECT().CheckProductGrants(user, opts.Product, auth.ActDeleteProcess).Return(nil)
	s.processRepo.EXPECT().GetByID(ctx, opts.Product, processID).Return(nil, nil)
	s.processRepo.EXPECT().Delete(ctx, opts.Product, processID).Return(nil)

	returnedProcessID, err := s.processService.DeleteProcess(ctx, user, opts)
	s.Require().NoError(err)

	s.Equal(processID, returnedProcessID)
}
