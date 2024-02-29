//go:build unit

package version_test

import (
	"context"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestGetLatest_OK() {
	// GIVEN
	var (
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		ctx         = context.Background()
		testVersion = testhelpers.NewVersionBuilder().WithTag("test-tag").Build()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.versionRepo.EXPECT().GetLatest(ctx, productID).Return(testVersion, nil)

	actual, err := s.handler.GetLatest(ctx, user, productID)
	s.Require().NoError(err)

	s.Equal(testVersion, actual)
}

func (s *versionSuite) TestGetLatest_ProductNotFound() {
	// GIVEN
	var (
		user      = testhelpers.NewUserBuilder().Build()
		productID = "product-1"
		ctx       = context.Background()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)

	// WHEN
	_, err := s.handler.GetLatest(ctx, user, productID)

	// THEN
	s.ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *versionSuite) TestGetLatest_VersionNotFoundWithExistingProduct() {
	// GIVEN
	var (
		user      = testhelpers.NewUserBuilder().Build()
		productID = "product-1"
		ctx       = context.Background()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.versionRepo.EXPECT().GetLatest(ctx, productID).Return(nil, version.ErrProductExistsNoVersions)

	// WHEN
	_, err := s.handler.GetLatest(ctx, user, productID)

	// THEN
	s.ErrorIs(err, version.ErrProductExistsNoVersions)
}
