//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestSearchByProduct_WithFilter_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	filter := &repository.ListVersionsFilter{
		Status: entity.VersionStatusCreated,
	}

	s.versionRepo.EXPECT().SearchByProduct(ctx, product.ID, filter).Return(nil, nil)

	_, err := s.handler.SearchByProduct(ctx, user, product.ID, filter)
	s.NoError(err)
}

func (s *versionSuite) TestSearchByProduct_InvalidFilter_Error() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	filter := &repository.ListVersionsFilter{
		Status: "INVALID",
	}

	_, err := s.handler.SearchByProduct(ctx, user, product.ID, filter)
	s.ErrorIs(err, entity.ErrInvalidVersionStatus)
}

func (s *versionSuite) TestSearchByProduct_Error() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	testError := errors.New("error")

	s.versionRepo.EXPECT().SearchByProduct(ctx, product.ID, nil).Return(nil, testError)

	_, err := s.handler.SearchByProduct(ctx, user, product.ID, nil)
	s.ErrorIs(err, testError)
}
