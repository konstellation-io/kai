//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/repository"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestListByProduct_NoFilter_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	s.versionRepo.EXPECT().ListVersionsByProduct(ctx, product.ID, nil).Return(nil, nil)

	_, err := s.handler.ListVersionsByProduct(ctx, user, product.ID, nil)
	s.NoError(err)
}

func (s *versionSuite) TestListByProduct_WithFilter_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	filter := &repository.ListVersionsFilter{
		Status: "CREATED",
	}

	s.versionRepo.EXPECT().ListVersionsByProduct(ctx, product.ID, filter).Return(nil, nil)

	_, err := s.handler.ListVersionsByProduct(ctx, user, product.ID, filter)
	s.NoError(err)
}

func (s *versionSuite) TestListByProduct_InvalidFilter_Error() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	filter := &repository.ListVersionsFilter{
		Status: "INVALID",
	}

	_, err := s.handler.ListVersionsByProduct(ctx, user, product.ID, filter)
	s.ErrorIs(err, entity.ErrInvalidVersionStatus)
}

func (s *versionSuite) TestListByProduct_Error() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	testError := errors.New("error")

	s.versionRepo.EXPECT().ListVersionsByProduct(ctx, product.ID, nil).Return(nil, testError)

	_, err := s.handler.ListVersionsByProduct(ctx, user, product.ID, nil)
	s.ErrorIs(err, testError)
}
