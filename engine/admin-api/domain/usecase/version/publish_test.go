//go:build unit

package version_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/konstellation-io/kai/engine/admin-api/domain/entity"
	"github.com/konstellation-io/kai/engine/admin-api/domain/service/auth"
	"github.com/konstellation-io/kai/engine/admin-api/domain/usecase/version"
	"github.com/konstellation-io/kai/engine/admin-api/testhelpers"
)

func (s *versionSuite) TestPublish_OK() {
	// GIVEN a valid user and published version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusStarted).
		Build()
	product := testhelpers.NewProductBuilder().Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publishing the version
	publishingVersion, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN the version status is publishing
	s.Require().NoError(err)
	s.Assert().Equal(entity.VersionStatusPublishing, publishingVersion.Status)

	publishedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Assert().Equal(user.Email, *publishedVersion.PublicationAuthor)
	s.Assert().Equal(entity.VersionStatusPublished, publishedVersion.Status)
}

func (s *versionSuite) TestPublishing_ErrorUserNotAuthorized() {
	// GIVEN an unauthorized user and a started version
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := &entity.Version{Tag: versionTag}
	product := testhelpers.NewProductBuilder().Build()

	expectedError := errors.New("unauthorized")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(expectedError)

	// WHEN publishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: expectedVer.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionNotFound() {
	// GIVEN a valid user and a version not found
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	expectedVer := testhelpers.NewVersionBuilder().Build()
	product := testhelpers.NewProductBuilder().Build()

	expectedError := errors.New("version not found")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, expectedVer.Tag).Return(nil, expectedError)

	// WHEN unpublishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: expectedVer.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionCannotBePublished() {
	// GIVEN a valid user and a version that cannot be unpublished
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithTag(versionTag).
		WithStatus(entity.VersionStatusCreated).
		Build()
	product := testhelpers.NewProductBuilder().Build()

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	// WHEN unpublishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrVersionCannotBePublished)
}

func (s *versionSuite) TestPublish_ProductWithVersionAlreadyPublished() {
	// GIVEN a valid user and a product with a published version
	var (
		ctx         = context.Background()
		user        = testhelpers.NewUserBuilder().Build()
		testVersion = "test-version"

		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(testhelpers.StrPointer("another-version")).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)

	// WHEN publishing the version
	_, _, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: testVersion,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrProductAlreadyPublished)
}

func (s *versionSuite) TestPublish_ErrorPublishingVersion() {
	// GIVEN a valid user and a published version, but error during publishing
	var (
		ctx  = context.Background()
		user = testhelpers.NewUserBuilder().Build()
		vers = testhelpers.NewVersionBuilder().
			WithStatus(entity.VersionStatusStarted).
			Build()

		product = testhelpers.NewProductBuilder().Build()

		expectedError = errors.New("publish error in k8s service")
		errStrMatcher = newStringContainsMatcher(expectedError.Error())
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, expectedError)

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), product.ID, vers.Tag, errStrMatcher).Return(nil)

	// WHEN no error in the initial return (the versionService publish is executed if a goroutine)
	_, notifyCh, _ := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Equal(entity.VersionStatusError, failedVersion.Status)
}

func (s *versionSuite) TestPublish_ErrorInStatusAndRegisteringAction() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(nil).
		Build()

	expectedURLs := map[string]string{
		"test-trigger": "test-url",
	}

	expectedError := errors.New("error registering user activity")
	errStrMatcher := newStringContainsMatcher(expectedError.Error())

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(expectedURLs, nil)

	versionMatcher := newVersionMatcher(vers)

	s.productRepo.EXPECT().Update(gomock.Any(), product).Times(2).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, versionMatcher).Times(2).Return(nil)

	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, gomock.Any(), "publishing").
		Return(expectedError)

	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, versionMatcher).Return(nil)
	// WHEN publishing the version
	v, notifyCh, err := s.handler.Publish(ctx, user, version.PublishParams{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	s.versionRepo.EXPECT().SetErrorStatusWithError(gomock.Any(), product.ID, vers.Tag, errStrMatcher).Return(nil)

	// THEN no error is returned (error happens in goroutine)
	s.Require().NoError(err)
	s.Equal(entity.VersionStatusPublishing, v.Status)

	failedVersion, ok := <-notifyCh
	s.Require().True(ok)

	s.Assert().Equal(entity.VersionStatusError, failedVersion.Status)
}
