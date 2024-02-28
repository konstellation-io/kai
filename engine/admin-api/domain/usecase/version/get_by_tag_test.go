//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestGetByTag_OK() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().WithTag("test-tag").Build()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)

	// WHEN
	actual, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.Require().NoError(err)
	s.Equal(testVersion, actual)
}

func (s *versionSuite) TestGetByTag_ProductNotFound() {
	// GIVEN
	var (
		user      = testhelpers.NewUserBuilder().Build()
		productID = "product-1"
		ctx       = context.Background()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(nil, usecase.ErrProductNotFound)

	// WHEN
	_, err := s.handler.GetByTag(ctx, user, productID, "test-tag")

	// THEN
	s.ErrorIs(err, usecase.ErrProductNotFound)
}

func (s *versionSuite) TestGetByTag_Unauthorized() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().WithTag("test-tag").Build()
		expectedErr = errors.New("unauthorized")
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(expectedErr)

	// WHEN
	_, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.ErrorIs(err, expectedErr)
}

func (s *versionSuite) TestGetByTag_VersionNotFound() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().WithTag("test-tag").Build()
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(nil, version.ErrVersionNotFound)

	// WHEN
	_, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.ErrorIs(err, version.ErrVersionNotFound)
}

func (s *versionSuite) TestGetByTag_PublishedVersion_PublishedTriggers() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().
				WithTag("test-tag").
				WithStatus(entity.VersionStatusPublished).
				Build()
		expectedPublishedTriggers = []entity.PublishedTrigger{
			{Trigger: "test-trigger", URL: "test-url"},
		}
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)
	s.versionService.EXPECT().GetPublishedTriggers(ctx, productID).Return(expectedPublishedTriggers, nil)

	// WHEN
	actual, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.Require().NoError(err)
	s.Equal(expectedPublishedTriggers, actual.PublishedTriggers)
}

func (s *versionSuite) TestGetByTag_PublishedVersion_ErrorGettingPublishedTriggers() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().
				WithTag("test-tag").
				WithStatus(entity.VersionStatusPublished).
				Build()
		expectedErr = errors.New("error getting published triggers")
	)

	s.productRepo.EXPECT().GetByID(ctx, productID).Return(&entity.Product{}, nil)
	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)
	s.versionService.EXPECT().GetPublishedTriggers(ctx, productID).Return(nil, expectedErr)

	// WHEN
	_, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.ErrorIs(err, expectedErr)
}
