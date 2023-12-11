//go:build unit

package process_test

import (
	"context"
	"time"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *ProcessServiceTestSuite) TestListByProduct_WithTypeFilter() {
	var (
		ctx               = context.Background()
		filter            = repository.SearchFilter{ProcessType: entity.ProcessTypeTrigger}
		productProcesses  = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder(productID).Build()}
		kaiProcesses      = []*entity.RegisteredProcess{testhelpers.NewRegisteredProcessBuilder("kai").Build()}
		expectedProcesses = append(productProcesses, kaiProcesses...)
	)

	s.processRepo.EXPECT().SearchByProduct(ctx, productID, filter).Return(productProcesses, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(kaiProcesses, nil)

	returnedRegisteredProcess, err := s.processService.Search(ctx, user, productID, filter.ProcessType.String())
	s.Require().NoError(err)

	s.Equal(expectedProcesses, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_NoTypeFilter() {
	ctx := context.Background()

	filter := repository.SearchFilter{}
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

	s.processRepo.EXPECT().SearchByProduct(ctx, productID, filter).Return(expectedRegisteredProcess, nil)
	s.processRepo.EXPECT().GlobalSearch(ctx, filter).Return(nil, nil)

	returnedRegisteredProcess, err := s.processService.Search(ctx, user, productID, "")
	s.Require().NoError(err)

	s.Equal(expectedRegisteredProcess, returnedRegisteredProcess)
}

func (s *ProcessServiceTestSuite) TestListByProduct_InvalidTypeFilterFilter() {
	ctx := context.Background()

	typeFilter := "invalid type"

	_, err := s.processService.Search(ctx, user, productID, typeFilter)
	s.Require().Error(err)
}
