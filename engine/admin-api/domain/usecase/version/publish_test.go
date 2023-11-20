//go:build unit

package version_test

import (
	"context"
	"errors"
	"sync"
	"time"

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

	s.versionService.EXPECT().Publish(ctx, product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publishing the version
	urls, err := s.handler.Publish(ctx, user, product.ID, versionTag, "publishing")

	// THEN the version status is published and the published triggers urls are returned
	s.Require().NoError(err)
	s.Assert().Equal(expectedURLs, urls)

	s.Assert().Equal(user.Email, *vers.PublicationAuthor)
	s.Assert().Equal(entity.VersionStatusPublished, vers.Status)
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
	_, err := s.handler.Publish(ctx, user, product.ID, expectedVer.Tag, "publishing")

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
	_, err := s.handler.Publish(ctx, user, product.ID, expectedVer.Tag, "publishing")

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
	_, err := s.handler.Publish(ctx, user, product.ID, versionTag, "publishing")

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
	_, err := s.handler.Publish(ctx, user, product.ID, testVersion, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrProductAlreadyPublished)
}

func (s *versionSuite) TestPublish_ErrorPublishingVersion() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	product := testhelpers.NewProductBuilder().Build()

	expectedError := errors.New("publish error in k8s service")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)

	s.versionService.EXPECT().Publish(ctx, product.ID, vers.Tag).Return(nil, expectedError)

	// WHEN publishing the version
	_, err := s.handler.Publish(ctx, user, product.ID, vers.Tag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(ctx, product.ID, vers.Tag).Return(expectedURLs, nil)

	versionMatcher := newVersionMatcher(vers)

	s.productRepo.EXPECT().Update(gomock.Any(), product).Times(2).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, versionMatcher).Times(2).Return(nil)

	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, gomock.Any(), "publishing").
		Return(expectedError)

	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, versionMatcher).
		DoAndReturn(func(_, _, _ interface{}) error {
			wg.Done()
			return nil
		})
	// WHEN publishing the version
	urls, err := s.handler.Publish(ctx, user, product.ID, vers.Tag, "publishing")

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
	s.Assert().Nil(urls)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))

	s.Assert().Equal(entity.VersionStatusStarted, vers.Status)
}
