//go:build unit

package process_test

import (
	"context"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessHandlerTestSuite) TestListByProduct_WithTypeFilter() {
	var (
		ctx               = context.Background()
		filter            = repository.SearchFilter{}
		productProcesses  = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder(productID).Build()}
		kaiProcesses      = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder("kai").Build()}
		expectedProcesses = append(productProcesses, kaiProcesses...)
	)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, nil)
	s.processRepo.EXPECT().SearchByProduct(ctx, productID, &filter).Return(productProcesses, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, &filter).Return(kaiProcesses, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, productID, &filter)
	s.Require().NoError(err)

	s.Equal(expectedProcesses, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestListByProduct_NoTypeFilter() {
	ctx := context.Background()

	expectedRegisteredProcess := []*entity.RegisteredProcess{
		{
			ID:         "test-id",
			Name:       processName,
			Version:    version,
			Type:       processType,
			Image:      "image",
			UploadDate: time.Now(),
			Owner:      userID,
		},
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, nil)
	s.processRepo.EXPECT().SearchByProduct(ctx, productID, nil).Return(expectedRegisteredProcess, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, nil).Return(nil, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, productID, nil)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestListByProduct_InvalidTypeFilterFilter() {
	ctx := context.Background()

	filter := repository.SearchFilter{
		ProcessType: "invalid",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, nil)

	_, err := s.processHandler.Search(ctx, user, productID, &filter)
	s.Require().Error(err)
}

func (s *ProcessHandlerTestSuite) TestListByProduct_ProductDoesNotExist() {
	ctx := context.Background()

	filter := repository.SearchFilter{
		ProcessType: "invalid",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)

	_, err := s.processHandler.Search(ctx, user, productID, &filter)
	s.Require().ErrorIs(err, usecase.ErrProductNotFound)
}
