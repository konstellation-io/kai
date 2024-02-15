//go:build unit

package process_test

import (
	"context"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessHandlerTestSuite) TestListByProduct_WithTypeFilter() {
	var (
		ctx               = context.Background()
		filter            = repository.SearchFilter{ProcessType: entity.ProcessTypeTrigger}
		productProcesses  = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder(_productID).Build()}
		kaiProcesses      = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder("kai").Build()}
		expectedProcesses = append(productProcesses, kaiProcesses...)
	)

	s.processRepo.EXPECT().SearchByProduct(ctx, _productID, filter).Return(productProcesses, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(kaiProcesses, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, _productID, filter.ProcessType.String())
	s.Require().NoError(err)

	s.Equal(expectedProcesses, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestListByProduct_NoTypeFilter() {
	ctx := context.Background()

	filter := repository.SearchFilter{}
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

	s.processRepo.EXPECT().SearchByProduct(ctx, _productID, filter).Return(expectedRegisteredProcess, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(nil, nil)

	returnedRegisteredProcess, err := s.processHandler.Search(ctx, user, _productID, "")
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessHandlerTestSuite) TestListByProduct_InvalidTypeFilterFilter() {
	ctx := context.Background()

	typeFilter := "invalid type"

	_, err := s.processHandler.Search(ctx, user, _productID, typeFilter)
	s.Require().Error(err)
}
