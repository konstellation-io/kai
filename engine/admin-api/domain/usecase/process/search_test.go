//go:build unit

package process_test

import (
	"context"
	"errors"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessHandlerTestSuite) TestSearch_WithTypeFilter() {
	var (
		ctx               = context.Background()
		filter            = repository.SearchFilter{}
		productProcesses  = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder(_productID).Build()}
		kaiProcesses      = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder("kai").Build()}
		expectedProcesses = append(productProcesses, kaiProcesses...)
	)

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(nil, nil)
	s.processRepo.EXPECT().SearchByProduct(ctx, _productID, &filter).Return(productProcesses, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, &filter).Return(kaiProcesses, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, _productID, &filter)
	s.Require().NoError(err)

	s.Equal(expectedProcesses, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestSearch_NoTypeFilter() {
	ctx := context.Background()

	expectedRegisteredProcess := []*entity.RegisteredProcess{
		{
			ID:         "test-id",
			Name:       _processName,
			Version:    _version,
			Type:       _processType,
			Image:      "image",
			UploadDate: time.Now(),
			Owner:      _userID,
		},
	}

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(nil, nil)
	s.processRepo.EXPECT().SearchByProduct(ctx, _productID, nil).Return(expectedRegisteredProcess, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, nil).Return(nil, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, _productID, nil)
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestSearch_InvalidTypeFilterFilter() {
	ctx := context.Background()

	filter := &repository.SearchFilter{
		ProcessType: "invalid",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(nil, nil)

	_, err := s.processHandler.Search(ctx, user, _productID, filter)
	s.Require().Error(err)
}

func (s *ProcessHandlerTestSuite) TestSearch_ProductDoesNotExist() {
	ctx := context.Background()

	filter := repository.SearchFilter{
		ProcessType: "invalid",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActViewRegisteredProcesses).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, _productID).Return(nil, usecase.ErrProductNotFound)

	_, err := s.processHandler.Search(ctx, user, _productID, &filter)
	s.Require().ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *ProcessHandlerTestSuite) TestSearch_Unauthorized() {
	ctx := context.Background()
	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, _productID, auth.ActViewRegisteredProcesses).Return(expectedError)

	_, err := s.processHandler.Search(ctx, user, _productID, &repository.SearchFilter{})
	s.Require().ErrorIs(err, expectedError)
}
