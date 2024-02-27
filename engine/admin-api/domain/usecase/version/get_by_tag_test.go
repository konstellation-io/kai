//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestGetByTag() {
	// GIVEN
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		productID   = "product-1"
		testVersion = testhelpers.NewVersionBuilder().WithTag("test-tag").Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)

	// WHEN
	actual, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.Require().NoError(err)
	s.Equal(testVersion, actual)
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

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(expectedErr)

	// WHEN
	_, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.ErrorIs(err, expectedErr)
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

	s.accessControl.EXPECT().CheckProductGrants(user, productID, auth.ActViewProduct).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, productID, testVersion.Tag).Return(testVersion, nil)
	s.versionService.EXPECT().GetPublishedTriggers(ctx, productID).Return(expectedPublishedTriggers, nil)

	// WHEN
	actual, err := s.handler.GetByTag(ctx, user, productID, testVersion.Tag)

	// THEN
	s.Require().NoError(err)
	s.Equal(expectedPublishedTriggers, actual.PublishedTriggers)
}
