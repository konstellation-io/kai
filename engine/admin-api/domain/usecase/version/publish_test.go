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

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publishing the version
	actualURLs, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN the version status is publishing
	s.Require().NoError(err)
	s.Equal(expectedURLs, actualURLs)

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
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
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
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: expectedVer.Tag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, expectedError)
}

func (s *versionSuite) TestPublish_ErrorVersionNotStarted_NoForce() {
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
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Assert().ErrorIs(err, version.ErrVersionIsNotStarted)
}

func (s *versionSuite) TestPublish_ProductWithVersionAlreadyPublished_NoForce() {
	// GIVEN a valid user and a product with a published version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(testhelpers.StrPointer("another-version")).
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusStarted).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	// WHEN publishing the version
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
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
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, expectedError)

	// WHEN no error in the initial return (the versionService publish is executed if a goroutine)
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	s.Error(err)
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	expectedError := errors.New("error registering user activity")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, nil)

	versionMatcher := newVersionMatcher(vers)

	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, versionMatcher).Return(nil)

	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, gomock.Any(), "publishing").
		Return(expectedError)

	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, versionMatcher).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).DoAndReturn(func(_, _ any) error {
		wg.Done()
		return nil
	})

	// WHEN
	urls, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	// THEN
	s.Assert().ErrorIs(err, expectedError)
	s.Assert().Nil(urls)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))

	s.Assert().Equal(entity.VersionStatusStarted, vers.Status)
}

func (s *versionSuite) TestPublish_AnotherVersionPublished_Forced() {
	// GIVEN a valid user and a non started version
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		oldVersion = testhelpers.NewVersionBuilder().
				WithTag("old-version").
				WithStatus(entity.VersionStatusPublished).
				Build()
		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(&oldVersion.Tag).
			Build()

		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusStarted).
			Build()
		expectedURLs = map[string]string{
			"test-trigger": "test-url",
		}
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.versionRepo.EXPECT().SetStatus(ctx, product.ID, oldVersion.Tag, entity.VersionStatusStarted).Return(nil)

	s.versionService.EXPECT().Publish(ctx, product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)

	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(nil)

	// WHEN publish the version with the param Force set to true
	actualURLs, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
		Force:      true,
	})

	// THEN an error is returned
	s.Require().NoError(err)
	s.Equal(expectedURLs, actualURLs)
	s.Equal(entity.VersionStatusPublished, vers.Status)
}

func (s *versionSuite) TestPublish_AnotherVersionPublished_Forced_RegisterActionError() {
	// GIVEN a valid user and a non started version
	var (
		ctx        = context.Background()
		user       = testhelpers.NewUserBuilder().Build()
		oldVersion = testhelpers.NewVersionBuilder().
				WithTag("old-version").
				WithStatus(entity.VersionStatusPublished).
				Build()
		product = testhelpers.NewProductBuilder().
			WithPublishedVersion(&oldVersion.Tag).
			Build()

		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusStarted).
			Build()
		expectedURLs = map[string]string{
			"test-trigger": "test-url",
		}

		wg = sync.WaitGroup{}
	)

	wg.Add(1)

	expecterError := errors.New("register action error")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	s.versionRepo.EXPECT().SetStatus(ctx, product.ID, oldVersion.Tag, entity.VersionStatusStarted).Return(nil)

	s.versionService.EXPECT().Publish(ctx, product.ID, vers.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)

	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, vers, "publishing").Return(expecterError)

	s.versionService.EXPECT().Publish(ctx, product.ID, oldVersion.Tag).Return(expectedURLs, nil)
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.versionRepo.EXPECT().SetStatus(gomock.Any(), product.ID, oldVersion.Tag, entity.VersionStatusPublished).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).DoAndReturn(func(_, _ any) error {
		wg.Done()
		return nil
	})

	// WHEN publish the version with the param Force set to true
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
		Force:      true,
	})

	// THEN an error is returned
	s.Require().ErrorIs(err, expecterError)
	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, _waitGroupTimeout))
}

func (s *versionSuite) TestPublish_FailsIfTheVersionIsAlreadyPublished() {
	// GIVEN a valid user and a non started version
	var (
		ctx     = context.Background()
		user    = testhelpers.NewUserBuilder().Build()
		product = testhelpers.NewProductBuilder().
			Build()
		vers = testhelpers.NewVersionBuilder().
			WithTag(versionTag).
			WithStatus(entity.VersionStatusPublished).
			Build()
	)

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, versionTag).Return(vers, nil)

	// WHEN publish the version with the param Force set to true
	_, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: versionTag,
		Comment:    "publishing",
	})

	// THEN an error is returned
	s.Require().ErrorIs(err, version.ErrVersionAlreadyPublished)
}

func (s *versionSuite) TestPublish_ErrorInStatusAndRegisteringAction_CompensationError() {
	// GIVEN a valid user and a published version, but error during publishing
	ctx := context.Background()
	user := testhelpers.NewUserBuilder().Build()
	vers := testhelpers.NewVersionBuilder().
		WithStatus(entity.VersionStatusStarted).
		Build()

	product := testhelpers.NewProductBuilder().
		WithPublishedVersion(nil).
		Build()

	wg := sync.WaitGroup{}
	wg.Add(1)

	expectedError := errors.New("error registering user activity")

	s.accessControl.EXPECT().CheckProductGrants(user, product.ID, auth.ActPublishVersion).Return(nil)
	s.productRepo.EXPECT().GetByID(ctx, product.ID).Return(product, nil)
	s.versionRepo.EXPECT().GetByTag(ctx, product.ID, vers.Tag).Return(vers, nil)

	s.versionService.EXPECT().Publish(gomock.Any(), product.ID, vers.Tag).Return(nil, nil)

	versionMatcher := newVersionMatcher(vers)

	s.productRepo.EXPECT().Update(ctx, product).Return(nil)
	s.versionRepo.EXPECT().Update(product.ID, versionMatcher).Return(nil)

	s.userActivityInteractor.EXPECT().RegisterPublishAction(user.Email, product.ID, gomock.Any(), "publishing").
		Return(expectedError)

	// Compensations
	s.versionService.EXPECT().Unpublish(gomock.Any(), product.ID, vers).Return(errors.New("unpublish error"))
	s.versionRepo.EXPECT().Update(product.ID, vers).Return(nil)
	s.productRepo.EXPECT().Update(gomock.Any(), product).Return(nil)
	s.versionRepo.EXPECT().SetCriticalStatusWithError(gomock.Any(), product.ID, vers.Tag, "unpublish error").
		DoAndReturn(func(_, _, _, _ any) error {
			wg.Done()

			return nil
		})

	// WHEN
	urls, err := s.handler.Publish(ctx, user, version.PublishOpts{
		ProductID:  product.ID,
		VersionTag: vers.Tag,
		Comment:    "publishing",
	})

	// THEN
	s.Assert().ErrorIs(err, expectedError)
	s.Assert().Nil(urls)

	s.Require().NoError(testhelpers.WaitOrTimeout(&wg, 1*time.Second))
}
